set WS=CreateObject("WScript.Shell")
for each ps in GetObject("winmgmts:\\.\root\cimv2:win32_process").instances_
if ps.name ="jane.exe" then
    WS.popup "程序已经在运行了",1,"提示",4144
    WScript.quit
end if

WS.Run("cmd /c jane.exe",0)
WS.popup "jane.exe 程序启动完毕",1,"提示",4144
