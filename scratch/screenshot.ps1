Add-Type -AssemblyName System.Windows.Forms
Add-Type -AssemblyName System.Drawing
\ = [System.Windows.Forms.Screen]::PrimaryScreen.Bounds
\ = New-Object System.Drawing.Bitmap \.Width, \.Height
\ = [System.Drawing.Graphics]::FromImage(\)
\.CopyFromScreen(0, 0, 0, 0, \.Size)
\.Save("scratch/screen.png", [System.Drawing.Imaging.ImageFormat]::Png)
\.Dispose()
\.Dispose()
