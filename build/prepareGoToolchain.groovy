import org.apache.commons.compress.archivers.tar.TarArchiveEntry
import org.apache.commons.compress.archivers.tar.TarArchiveInputStream
import org.apache.commons.compress.archivers.zip.ZipArchiveInputStream
import org.apache.commons.io.FileUtils
import org.apache.commons.io.IOUtils
import org.apache.commons.lang.StringUtils
import org.apache.commons.lang.SystemUtils
import org.apache.maven.project.MavenProject
import org.slf4j.LoggerFactory

import java.util.regex.Pattern
import java.util.zip.GZIPInputStream

public class PrepareGoToolchain {

    private final static LOG = LoggerFactory.getLogger("prepareGoToolchain.groovy")
    private final static REMOVE_LEADING_GO_PATH_PATTERN = Pattern.compile("^(|\\./)go/")
    private final static PLATFORM_PATTERN = Pattern.compile("\\s*([a-z0-9]+)-([a-z0-9]+)\\s*", Pattern.MULTILINE)

    private static URL determinateDownloadUrl(String goVersion) {
        final String suffix;
        if (SystemUtils.IS_OS_WINDOWS) {
            suffix = "windows-amd64.zip"
        } else if (SystemUtils.IS_OS_LINUX) {
            suffix = "linux-amd64.tar.gz"
        } else if (SystemUtils.IS_OS_MAC_OSX) {
            suffix = "darwin-amd64.tar.gz"
        } else {
            throw new IllegalArgumentException("The current operating system is not supported by this build system: ${SystemUtils.OS_ARCH} ${SystemUtils.OS_ARCH}")
        }
        return new URL("https://storage.googleapis.com/golang/go${goVersion}.${suffix}")
    }

    private static void unTarGz(InputStream is, File target) {
        final gzip = new GZIPInputStream(is)
        final archive = new TarArchiveInputStream(gzip)
        def TarArchiveEntry entry = archive.nextTarEntry
        while (entry != null) {
            final entryFile = new File(target, REMOVE_LEADING_GO_PATH_PATTERN.matcher(entry.name).replaceFirst("")).canonicalFile
            if (entry.directory) {
                FileUtils.forceMkdir(entryFile)
            } else {
                LOG.debug("Write: {}...", entryFile)
                FileUtils.forceMkdir(entryFile.parentFile)
                final os = new FileOutputStream(entryFile)
                try {
                    IOUtils.copy(archive, os)
                } finally {
                    IOUtils.closeQuietly(os)
                }
                entryFile.setExecutable(
                        (entry.mode | 0100) > 0,
                        (entry.mode | 0001) == 0
                )
            }
            entry = archive.nextTarEntry
        }
    }

    private static void unZip(InputStream is, File target) {
        final archive = new ZipArchiveInputStream(new BufferedInputStream(is), "utf-8", false, true)
        def entry = archive.nextZipEntry
        while (entry != null) {
            final entryFile = new File(target, REMOVE_LEADING_GO_PATH_PATTERN.matcher(entry.name).replaceFirst("")).canonicalFile
            if (entry.directory) {
                FileUtils.forceMkdir(entryFile)
            } else {
                LOG.debug("Write: {}...", entryFile)
                FileUtils.forceMkdir(entryFile.getParentFile())
                final os = new FileOutputStream(entryFile)
                try {
                    IOUtils.copy(archive, os)
                } finally {
                    IOUtils.closeQuietly(os)
                }
            }
            entry = archive.nextZipEntry
        }
    }

    private static void downloadAndExtract(String goVersion, String targetPath, boolean force) {
        final downloadUrl = determinateDownloadUrl(goVersion)
        final target = new File(targetPath)

        final targetMarker = new File("${target.path}/.downloaded")
        def previousDownloadLocation = ""
        try {
            previousDownloadLocation = FileUtils.readFileToString(targetMarker)
        } catch (IOException ignored) {}
        if (force || !downloadUrl.toExternalForm().equals(previousDownloadLocation)) {
            FileUtils.deleteDirectory(target)
            FileUtils.forceMkdir(target)

            LOG.info("Going to download {} and extract it to {}...", downloadUrl, target)
            final is = downloadUrl.openStream()
            try {
                if (downloadUrl.toExternalForm().endsWith(".tar.gz")) {
                    unTarGz(is, target)
                } else if (downloadUrl.toExternalForm().endsWith(".zip")) {
                    unZip(is, target)
                } else {
                    throw new IllegalStateException("Does not support download archive of type ${downloadUrl.toExternalForm()}.")
                }
            } finally {
                IOUtils.closeQuietly(is)
            }
            new File(target, "bin/go").setExecutable(true, false)
            new File(target, "src/make.bash").setExecutable(true, false)
            LOG.info("Going to download {} and extract it to {}... DONE!", downloadUrl, target)
            FileUtils.writeStringToFile(targetMarker, downloadUrl.toExternalForm())
        } else {
            LOG.info("{} already downloaded to {}.", downloadUrl, target)
        }
    }

    private static void buildGoForAllOf(String platformsListing, String goroot, boolean force) {
        final matcher = PLATFORM_PATTERN.matcher(platformsListing)
        def i = 0
        while (matcher.find(i)) {
            buildGoFor(matcher.group(1), matcher.group(2), goroot, force)
            i = matcher.end()
        }
    }

    private static void buildGoFor(String os, String architecture, String goroot, boolean force) {
        final buildMarker = new File("${goroot}/pkg/${os}_${architecture}/.builded")
        if (force || ! buildMarker.file) {
            final sourceDirectory = new File("${goroot}/src")
            final makeScriptExtension = SystemUtils.IS_OS_WINDOWS ? "bat" : "bash"
            final makeScript = new File(sourceDirectory, "make.${makeScriptExtension}")
            final pb = new ProcessBuilder(makeScript.path, "--no-clean")
            pb.directory(sourceDirectory)
            pb.environment().put("GOROOT", goroot)
            pb.environment().put("GOROOT_BOOTSTRAP", goroot)
            pb.environment().put("GOOS", os)
            pb.environment().put("GOARCH", architecture)
            pb.environment().put("CGO_ENABLED", "0")
            pb.redirectOutput(ProcessBuilder.Redirect.INHERIT)
            pb.redirectErrorStream(true)

            LOG.info("Going to build go toolchain for {}-{}...", os, architecture)
            final process = pb.start()
            final exitCode = process.waitFor()
            if (exitCode != 0) {
                throw new RuntimeException("Command ${pb.command()} failed with exitCode ${exitCode}.")
            }
            FileUtils.writeStringToFile(buildMarker, "")
            LOG.info("Going to build go toolchain for {}-{}... DONE!", os, architecture)
        } else {
            LOG.info("go toolchain for {}-{} already build.", os, architecture)
        }
    }

    public static run(MavenProject project){
        final goVersion = project.properties.getProperty('project.versions.go')
        if (StringUtils.isEmpty(project.properties.getProperty("project.go.root"))) {
            project.properties.setProperty("project.go.root", "${System.getProperty("user.home", project.basedir.path)}/.go-bootstrap/${goVersion}")
        }

        if (StringUtils.isEmpty(project.properties.getProperty("project.go.path"))) {
            project.properties.setProperty("project.go.path", "${project.build.directory}/gopath")
        }

        downloadAndExtract(
                goVersion,
                project.properties.getProperty('project.go.root'),
                Boolean.TRUE.toString().equalsIgnoreCase(project.properties.getProperty('project.go.forceDownloadToolchain'))
        )

        buildGoForAllOf(
                project.properties.getProperty('project.build.targetPlatforms'),
                project.properties.getProperty('project.go.root'),
                Boolean.TRUE.toString().equalsIgnoreCase(project.properties.getProperty('project.go.forceBuildToolchain'))
        )
    }
}

//noinspection GrUnresolvedAccess,GroovyAssignabilityCheck
PrepareGoToolchain.run(project)


