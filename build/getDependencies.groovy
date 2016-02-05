import org.apache.commons.io.FileUtils
import org.apache.commons.lang.SystemUtils
import org.apache.maven.project.MavenProject
import org.slf4j.LoggerFactory

import java.util.regex.Pattern

public class GetDependencies {

    private final static LOG = LoggerFactory.getLogger("getDependencies.groovy")
    private final static PLATFORM_PATTERN = Pattern.compile("\\s*([a-z0-9]+)-([a-z0-9]+)\\s*", Pattern.MULTILINE)

    private static void getDependencies(String goroot, String gopath, String dependenciesAsString) {
        if (dependenciesAsString != null && !dependenciesAsString.isEmpty()) {
            final dependencies = dependenciesAsString.split("[\\s,;]")
            for (final plainDependency : dependencies) {
                final dependency = plainDependency.trim()
                if (! dependency.isEmpty()) {
                    getDependency(goroot, gopath, dependency)
                }
            }
        }
    }

    private static void getDependency(String goroot, String gopath, String dependency) {
        final gopathFile = new File(gopath)
        FileUtils.forceMkdir(gopathFile)
        final binarySuffix = SystemUtils.IS_OS_WINDOWS ? ".exe" : ""
        final pb = new ProcessBuilder("${goroot}/bin/go${binarySuffix}", "get", dependency)
        pb.directory(gopathFile)
        pb.environment().put("GOROOT", goroot)
        pb.environment().put("GOPATH", gopath)
        pb.redirectOutput(ProcessBuilder.Redirect.INHERIT)
        pb.redirectErrorStream(true)

        LOG.info("Going to get http://{}...", dependency)
        final process = pb.start()
        final exitCode = process.waitFor()
        if (exitCode != 0) {
            throw new RuntimeException("Command ${pb.command()} failed with exitCode ${exitCode}.")
        }
        LOG.info("Going to get http://{}... DONE!", dependency)
    }

    public static run(MavenProject project, String dependenciesAsString) {
        getDependencies(
                project.properties.getProperty('project.go.root'),
                project.properties.getProperty('project.go.path'),
                dependenciesAsString
        )
    }
}

//noinspection GrUnresolvedAccess,GroovyAssignabilityCheck
GetDependencies.run(project, properties['dependencies'])


