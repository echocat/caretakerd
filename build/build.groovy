import org.apache.commons.io.FileUtils
import org.apache.commons.lang.SystemUtils
import org.apache.maven.project.MavenProject
import org.slf4j.LoggerFactory

import java.util.regex.Pattern

public class Build {

    private final static LOG = LoggerFactory.getLogger("build.groovy")
    private final static PLATFORM_PATTERN = Pattern.compile("\\s*([a-z0-9]+)-([a-z0-9]+)\\s*", Pattern.MULTILINE)

    private static String buildForAllOf(String targetPlatforms, String goroot, String gopath, String pkg, String output) {
        def result = ""
        if (targetPlatforms == null || targetPlatforms.trim().isEmpty()) {
            result = build(null, null, goroot, gopath, pkg, output)
        } else {
            final matcher = PLATFORM_PATTERN.matcher(targetPlatforms)
            def i = 0
            while (matcher.find(i)) {
                if (!result.isEmpty()) {
                    result += File.pathSeparator
                }
                result += build(matcher.group(1), matcher.group(2), goroot, gopath, pkg, output)
                i = matcher.end()
            }
        }
        return result
    }

    private static String build(String os, String architecture, String goroot, String gopath, String pkg, String output) {
        final gopathFile = new File(gopath)
        FileUtils.forceMkdir(gopathFile)
        final binarySuffix = SystemUtils.IS_OS_WINDOWS ? ".exe" : ""
        final String targetBinary
        if (os == null || architecture == null) {
            targetBinary = "${output}${binarySuffix}"
        } else if (os == "windows") {
            targetBinary = "${output}-${os}-${architecture}.exe"
        } else {
            targetBinary = "${output}-${os}-${architecture}"
        }
        final pb = new ProcessBuilder("${goroot}/bin/go${binarySuffix}", "build", "-o", targetBinary, pkg)
        pb.directory(gopathFile)
        pb.environment().put("GOROOT", goroot)
        pb.environment().put("GOPATH", gopath)
        if (os != null) {
            pb.environment().put("GOOS", os)
        }
        if (architecture != null) {
            pb.environment().put("GOARCH", architecture)
        }
        pb.redirectOutput(ProcessBuilder.Redirect.INHERIT)
        pb.redirectErrorStream(true)

        LOG.info("Going to build {} to {}...", pkg, targetBinary, os, architecture)
        final process = pb.start()
        final exitCode = process.waitFor()
        if (exitCode != 0) {
            throw new RuntimeException("Command ${pb.command()} failed with exitCode ${exitCode}.")
        }
        LOG.info("Going to build {} to {}... DONE!", pkg, targetBinary, os, architecture)
        return targetBinary
    }

    public static run(MavenProject project, String pkg, String output, String targetPlatforms, String exportToProperty){
        final build = buildForAllOf(
                targetPlatforms,
                project.properties.getProperty('project.go.root'),
                project.properties.getProperty('project.go.path'),
                pkg,
                output
        )
        if (exportToProperty != null && !exportToProperty.isEmpty()) {
            project.properties.setProperty(exportToProperty, build)
        }
    }
}

//noinspection GrUnresolvedAccess,GroovyAssignabilityCheck
Build.run(project, properties['pkg'], properties['output'], properties['targetPlatforms'], properties['exportToProperty'])


