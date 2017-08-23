package notifier

import (
	"text/template"
)

var (
	mailTemplate = template.Must(template.New("mail").Parse(`
<html>
<head>
    <title>*|Searchlight Alert|*</title>
<style type="text/css">
@font-face {
font-family: 'Roboto'; font-style: normal; font-weight: 100; src: local('Roboto Thin'), local('Roboto-Thin'), url('https://fonts.gstatic.com/s/roboto/v15/Jzo62I39jc0gQRrbndN6nfesZW2xOQ-xsNqO47m55DA.ttf') format('truetype');
}
@font-face {
font-family: 'Roboto'; font-style: normal; font-weight: 300; src: local('Roboto Light'), local('Roboto-Light'), url('https://fonts.gstatic.com/s/roboto/v15/Hgo13k-tfSpn0qi1SFdUfaCWcynf_cDxXwCLxiixG1c.ttf') format('truetype');
}
@font-face {
font-family: 'Roboto'; font-style: normal; font-weight: 400; src: local('Roboto'), local('Roboto-Regular'), url('https://fonts.gstatic.com/s/roboto/v15/zN7GBFwfMP4uA6AR0HCoLQ.ttf') format('truetype');
}
@font-face {
font-family: 'Roboto'; font-style: normal; font-weight: 500; src: local('Roboto Medium'), local('Roboto-Medium'), url('https://fonts.gstatic.com/s/roboto/v15/RxZJdnzeo3R5zSexge8UUaCWcynf_cDxXwCLxiixG1c.ttf') format('truetype');
}
@font-face {
font-family: 'Roboto'; font-style: normal; font-weight: 700; src: local('Roboto Bold'), local('Roboto-Bold'), url('https://fonts.gstatic.com/s/roboto/v15/d-6IYplOFocCacKzxwXSOKCWcynf_cDxXwCLxiixG1c.ttf') format('truetype');
}
@font-face {
font-family: 'Roboto'; font-style: normal; font-weight: 900; src: local('Roboto Black'), local('Roboto-Black'), url('https://fonts.gstatic.com/s/roboto/v15/mnpfi9pxYH-Go5UiibESIqCWcynf_cDxXwCLxiixG1c.ttf') format('truetype');
}
@font-face {
font-family: 'Roboto'; font-style: italic; font-weight: 100; src: local('Roboto Thin Italic'), local('Roboto-ThinItalic'), url('https://fonts.gstatic.com/s/roboto/v15/12mE4jfMSBTmg-81EiS-YS3USBnSvpkopQaUR-2r7iU.ttf') format('truetype');
}
@font-face {
font-family: 'Roboto'; font-style: italic; font-weight: 300; src: local('Roboto Light Italic'), local('Roboto-LightItalic'), url('https://fonts.gstatic.com/s/roboto/v15/7m8l7TlFO-S3VkhHuR0at50EAVxt0G0biEntp43Qt6E.ttf') format('truetype');
}
@font-face {
font-family: 'Roboto'; font-style: italic; font-weight: 400; src: local('Roboto Italic'), local('Roboto-Italic'), url('https://fonts.gstatic.com/s/roboto/v15/W4wDsBUluyw0tK3tykhXEfesZW2xOQ-xsNqO47m55DA.ttf') format('truetype');
}
@font-face {
font-family: 'Roboto'; font-style: italic; font-weight: 500; src: local('Roboto Medium Italic'), local('Roboto-MediumItalic'), url('https://fonts.gstatic.com/s/roboto/v15/OLffGBTaF0XFOW1gnuHF0Z0EAVxt0G0biEntp43Qt6E.ttf') format('truetype');
}
@font-face {
font-family: 'Roboto'; font-style: italic; font-weight: 700; src: local('Roboto Bold Italic'), local('Roboto-BoldItalic'), url('https://fonts.gstatic.com/s/roboto/v15/t6Nd4cfPRhZP44Q5QAjcC50EAVxt0G0biEntp43Qt6E.ttf') format('truetype');
}
@font-face {
font-family: 'Roboto'; font-style: italic; font-weight: 900; src: local('Roboto Black Italic'), local('Roboto-BlackItalic'), url('https://fonts.gstatic.com/s/roboto/v15/bmC0pGMXrhphrZJmniIZpZ0EAVxt0G0biEntp43Qt6E.ttf') format('truetype');
}
body {
background: #f2f2f2; color: #263249; font-family: Roboto, sans-serif; font-size: 18px; font-weight: 400; margin-top: 30px; padding: 0;
}
</style>
</head>
<body style="color: #263249; font-family: Roboto, sans-serif; font-size: 18px; font-weight: 400; margin-top: 30px; background: #f2f2f2; padding: 0;">
<div class="wrapper">
    <div class="email-body" style="border-radius: 5px; height: auto; max-width: 500px; overflow: hidden; background: #fff; margin: 0 auto; border: 1px solid #e9e9e9;">
        <div class="header" style="border-bottom-width: 1px; border-bottom-color: #e9e9e9; border-bottom-style: solid; background: #fbfbfb;">
            <div class="header_top">
                <div class="logo" style="text-align: center; padding: 40px 30px;" align="center">
                    <a href="http://www.appscode.com/" target="_blank" style="display: block;">
                        <img align="center" alt="" src="https://cdn.appscode.com/images/logo.png" style="width: 50%;" />
                    </a>
                </div>
                <!-- end logo -->
            </div>
            <!-- end header_top -->
        </div>
        <!-- end header -->
        <div class="body_content">
            <div class="content-top" style="width: 80%; margin: 0 auto; padding: 40px 0;">
                <!-- <h3>Mail alert</h3> -->
                <div class="notification">
                    <h4 style="font-size: 15px; font-weight: 400; margin: 0; padding: 4px 0;">This is an automated incident notification.</h4>
                    <p style="font-size: 14px; font-weight: normal; line-height: 25px; margin: 0 0 20px; padding: 0px;">Please see details below:</p>

                    <div class="cluster-info" style="margin-bottom: 20px; display: block;">
                        <h4 style="font-size: 15px; font-weight: 400; text-transform: capitalize; border-top-left-radius: 3px; border-top-right-radius: 3px; overflow: hidden; background: #fcfcfc; margin: 0; padding: 4px; border-color: #f0f0f0#f0f0f0#dddddd; border-style: solid; border-width: 1px 1px 0px;">Kubernetes Information</h4>
                        <table style="border-collapse: collapse; width: 100%; border: 0px solid #f0f0f0;">

                            {{ if .AlertNamespace }}
                            <tr>
                                <td style="text-align: left; vertical-align: top; color: #575757; font-size: 12px; padding: 6px; border: 1px solid #f0f0f0;" align="left" valign="top">Namespace</td>
                                <td style="text-align: left; vertical-align: top; color: #575757; font-size: 12px; padding: 6px; border: 1px solid #f0f0f0;" align="left" valign="top">{{ .AlertNamespace  }}</td>
                            </tr>
                            {{ end }}

                            {{ if .AlertType }}
                            <tr>
                                <td style="text-align: left; vertical-align: top; color: #575757; font-size: 12px; padding: 6px; border: 1px solid #f0f0f0;" align="left" valign="top">Alert Type</td>
                                <td style="text-align: left; vertical-align: top; color: #575757; font-size: 12px; padding: 6px; border: 1px solid #f0f0f0;" align="left" valign="top">{{ .AlertType  }}</td>
                            </tr>
                            {{ end }}

                            {{ if .AlertName }}
                            <tr>
                                <td style="text-align: left; vertical-align: top; color: #575757; font-size: 12px; padding: 6px; border: 1px solid #f0f0f0;" align="left" valign="top">Alert Name</td>
                                <td style="text-align: left; vertical-align: top; color: #575757; font-size: 12px; padding: 6px; border: 1px solid #f0f0f0;" align="left" valign="top">{{ .AlertName  }}</td>
                            </tr>
                            {{ end }}

                            {{ if .ObjectName }}
                            <tr>
                                <td style="text-align: left; vertical-align: top; color: #575757; font-size: 12px; padding: 6px; border: 1px solid #f0f0f0;" align="left" valign="top">Object Name</td>
                                <td style="text-align: left; vertical-align: top; color: #575757; font-size: 12px; padding: 6px; border: 1px solid #f0f0f0;" align="left" valign="top">{{ .ObjectName  }}</td>
                            </tr>
                            {{ end }}

                        </table>
                    </div>

                    <div class="cluster-info" style="margin-bottom: 20px; display: block;">
                        <h4 style="font-size: 15px; font-weight: 400; text-transform: capitalize; border-top-left-radius: 3px; border-top-right-radius: 3px; overflow: hidden; background: #fcfcfc; margin: 0; padding: 4px; border-color: #f0f0f0#f0f0f0#dddddd; border-style: solid; border-width: 1px 1px 0px;">Incident Event Information</h4>
                        <table style="border-collapse: collapse; width: 100%; border: 0px solid #f0f0f0;">

                            {{ if .IcingaHostName }}
                            <tr>
                                <td style="text-align: left; vertical-align: top; color: #575757; font-size: 12px; padding: 6px; border: 1px solid #f0f0f0;" align="left" valign="top">Host Name</td>
                                <td style="text-align: left; vertical-align: top; color: #575757; font-size: 12px; padding: 6px; border: 1px solid #f0f0f0;" align="left" valign="top">{{ .IcingaHostName }}</td>
                            </tr>
                            {{ end }}

                            {{ if .IcingaServiceName }}
                            <tr>
                                <td style="text-align: left; vertical-align: top; color: #575757; font-size: 12px; padding: 6px; border: 1px solid #f0f0f0;" align="left" valign="top">Service Name</td>
                                <td style="text-align: left; vertical-align: top; color: #575757; font-size: 12px; padding: 6px; border: 1px solid #f0f0f0;" align="left" valign="top">{{ .IcingaServiceName }}</td>
                            </tr>
                            {{ end }}

                            {{ if .IcingaCheckCommand }}
                            <tr>
                                <td style="text-align: left; vertical-align: top; color: #575757; font-size: 12px; padding: 6px; border: 1px solid #f0f0f0;" align="left" valign="top">Check Command</td>
                                <td style="text-align: left; vertical-align: top; color: #575757; font-size: 12px; padding: 6px; border: 1px solid #f0f0f0;" align="left" valign="top">{{ .IcingaCheckCommand }}</td>
                            </tr>
                            {{ end }}

                            {{ if .IcingaType }}
                            <tr>
                                <td style="text-align: left; vertical-align: top; color: #575757; font-size: 12px; padding: 6px; border: 1px solid #f0f0f0;" align="left" valign="top">Type</td>
                                <td style="font-weight: 600; font-size: 14px; text-align: left; vertical-align: top; color: #575757; background: #fafafa; padding: 6px; border: 1px solid #f0f0f0;" align="left" valign="top">{{ .IcingaType }}</td>
                            </tr>
                            {{ end }}

                            {{ if eq .IcingaState "OK" }}
                            <tr>
                                <td style="text-align: left; vertical-align: top; color: #575757; font-size: 12px; padding: 6px; border: 1px solid #f0f0f0;" align="left" valign="top">State</td>
                                <td style="text-align: left; vertical-align: top; color: #006400; font-size: 12px; padding: 6px; border: 1px solid #f0f0f0;" align="left" valign="top">{{ .IcingaState }}</td>
                            </tr>
                            {{ else if eq .IcingaState "CRITICAL" }}
                            <tr>
                                <td style="text-align: left; vertical-align: top; color: #575757; font-size: 12px; padding: 6px; border: 1px solid #f0f0f0;" align="left" valign="top">State</td>
                                <td style="text-align: left; vertical-align: top; color: #FF0000; font-size: 12px; padding: 6px; border: 1px solid #f0f0f0;" align="left" valign="top">{{ .IcingaState }}</td>
                            </tr>
                            {{ else if eq .IcingaState "WARNING" }}
                            <tr>
                                <td style="text-align: left; vertical-align: top; color: #575757; font-size: 12px; padding: 6px; border: 1px solid #f0f0f0;" align="left" valign="top">State</td>
                                <td style="text-align: left; vertical-align: top; color: #FF7F50; font-size: 12px; padding: 6px; border: 1px solid #f0f0f0;" align="left" valign="top">{{ .IcingaState }}</td>
                            </tr>
                            {{ else if .IcingaState }}
                            <tr>
                                <td style="text-align: left; vertical-align: top; color: #575757; font-size: 12px; padding: 6px; border: 1px solid #f0f0f0;" align="left" valign="top">State</td>
                                <td style="text-align: left; vertical-align: top; color: #575757; font-size: 12px; padding: 6px; border: 1px solid #f0f0f0;" align="left" valign="top">{{ .IcingaState }}</td>
                            </tr>
                            {{ end }}

                            {{ if .IcingaOutput }}
                            <tr>
                                <td style="text-align: left; vertical-align: top; color: #575757; font-size: 12px; padding: 6px; border: 1px solid #f0f0f0;" align="left" valign="top">Service Output</td>
                                <td><pre>{{ .IcingaOutput }}</pre></td>
                            </tr>
                            {{ end }}

                        </table>
                    </div>

                    <h2 style="font-size: 14px; font-weight: 400; margin: 10px 0 0; padding: 4px 0;"> Reported at <span style="border-bottom-width: 1px; border-bottom-color: #3bb778; border-bottom-style: dotted; margin: 0 8px;">{{ .IcingaTime }}</span></h2>
                </div>
            </div>
            <!-- end content-top -->
        </div>
        <!-- end body_content -->
        <div id="footer" style="padding: 30px 30px 20px; background: #fbfbfb; border-top-width: 1px; border-top-style: solid; border-top-color: #e9e9e9; color: #9b9b9b; font-size: 14px;">

            <div class="footer-right" style="display: block; margin-bottom: 20px; text-align: center; vertical-align: middle;" align="center">
                <ul style="margin: 0; padding: 0;">
                    <li style="border-radius: 50px; display: inline-block; height: 40px; margin-right: 10px; text-align: center; width: 40px; list-style: none;">
                        <a href="https://appscode.com/" target="_blank" style="color: #9b9b9b; display: inline; font-size: 20px; font-weight: normal; line-height: 25px; padding: 5px 10px;"><img align="center" alt="web" src="https://cdn.appscode.com/images/web.png"></a>
                    </li>
                    <li style="border-radius: 50px; display: inline-block; height: 40px; margin-right: 10px; text-align: center; width: 40px; list-style: none;">
                        <a href="https://twitter.com/AppsCodeHQ" target="_blank" style="color: #9b9b9b; display: inline; font-size: 20px; font-weight: normal; line-height: 25px; padding: 5px 10px;"><img align="center" alt="twitter" src="https://cdn.appscode.com/images/twitter.png"></a>
                    </li>
                    <li style="border-radius: 50px; display: inline-block; height: 40px; margin-right: 10px; text-align: center; width: 40px; list-style: none;">
                        <a href="https://www.facebook.com/appscode" target="_blank" style="color: #9b9b9b; display: inline; font-size: 20px; font-weight: normal; line-height: 25px; padding: 5px 10px;"><img align="center" alt="facebook" src="https://cdn.appscode.com/images/facebook.png"></a>
                    </li>

                    <li style="border-radius: 50px; display: inline-block; height: 40px; margin-right: 10px; text-align: center; width: 40px; list-style: none;">
                        <a href="https://github.com/appscode" target="_blank" style="color: #9b9b9b; display: inline; font-size: 20px; font-weight: normal; line-height: 25px; padding: 5px 10px;"><img align="center" alt="github" src="https://cdn.appscode.com/images/github.png"></a>
                    </li>
                    <li style="border-radius: 50px; display: inline-block; height: 40px; margin-right: 10px; text-align: center; width: 40px; list-style: none;">
                        <a href="https://hub.docker.com/u/appscode" target="_blank" style="color: #9b9b9b; display: inline; font-size: 20px; font-weight: normal; line-height: 25px; padding: 5px 10px;"><img align="center" alt="docker" src="https://cdn.appscode.com/images/docker.png"></a>
                    </li>
                </ul>
            </div>
            <div class="footer-left" style="margin-top: 30px; text-align: center;" align="center">
                <ul style="margin: 0; padding: 0;">
                    <li style="color: #aaa; display: block; font-size: 12px; font-weight: normal; line-height: 25px; list-style: none;">The Appscode Team</li>
                </ul>
            </div>
        </div>
    </div>
    <!-- end email-body -->
</div>
<!-- end wrapper -->
</body>
</html>
`))
)
