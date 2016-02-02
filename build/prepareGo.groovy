import org.apache.commons.compress.archivers.tar.TarArchiveEntry
import org.apache.commons.compress.archivers.tar.TarArchiveInputStream
import org.apache.commons.io.FileUtils
import org.apache.commons.io.IOUtils

import java.util.regex.Pattern
import java.util.zip.GZIPInputStream


project.properties.setProperty("GOROOT", "${project.build.directory}/go")
project.properties.setProperty("GOROOT_BOOTSTRAP", "${project.build.directory}/go_bootstrap")
project.properties.setProperty("GOPATH", "${project.build.directory}/gopath")

public void downloadAndExtract(String sourceUrl, String targetPath) {
    final removeLeadingGoPathPattern = Pattern.compile("^(|\\./)go/")

    final target = new File(targetPath)
    FileUtils.deleteDirectory(target)
    FileUtils.forceMkdir(target)

    final is = new URL(sourceUrl).openStream()
    try {
        final gzip = new GZIPInputStream(is)
        final tar = new TarArchiveInputStream(gzip)
        def TarArchiveEntry entry = null
        while ((entry = tar.nextTarEntry) != null) {
            final entryFile = new File(target, removeLeadingGoPathPattern.matcher(entry.name).replaceFirst(""))
            if (entry.directory) {
                FileUtils.forceMkdir(entryFile)
            } else {
                final os = new FileOutputStream(entryFile)
                try {
                    IOUtils.copy(tar, os)
                } finally {
                    IOUtils.closeQuietly(os)
                }
            }
        }
    } finally {
        IOUtils.closeQuietly(is)
    }
}

public static void buildGoFor(String os, String architecture) {
}

downloadAndExtract("https://storage.googleapis.com/golang/go${project.properties['GO_VERSION']}.linux-amd64.tar.gz", project.properties['GOROOT_BOOTSTRAP'])
downloadAndExtract("https://golang.org/dl/go${project.properties['GO_VERSION']}.src.tar.gz", project.properties['GOROOT'])

buildGoFor("windows", "amd64")