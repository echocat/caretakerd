The way how the binary reacts depends on the command/file name it was called. A simple rename or symlink will change the behavior of the binary.

| Pattern | Target command | Examples |
| ------- | ---------------| -------- |
| ``*/caretakerd*`` | [``caretakerd``](#commands.caretakerd) | ``/usr/bin/caretakerd``<br>``/usr/bin/caretakerd-linux-amd64``<br>``C:\Program Files\caretakerd\caretakerd-windows-amd64.exe`` |
| ``*/caretakerctl*`` | [``caretakerctl``](#commands.caretakerctl) | ``/usr/bin/caretakerctl``<br>``/usr/bin/caretakerctl-linux-amd64``<br>``C:\Program Files\caretakerd\caretakerctl-windows-amd64.exe`` |
| Everything else that does<br>not matches one of above. | [``caretaker``](#commands.caretaker) | ``/usr/bin/caretaker``<br>``/usr/bin/caretaker-linux-amd64``<br>``C:\Program Files\caretakerd\caretaker-windows-amd64.exe`` |
