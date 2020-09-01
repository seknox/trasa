using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Data;
using System.Drawing;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using System.Windows.Forms;

using System.Diagnostics;
using System.Runtime.InteropServices;

using Newtonsoft.Json;
using System.IO;
using System.Net;

struct trasaCPConfig
{
    public string trasaHost, serviceID, serviceKey, offlineUsers; public bool skipTLSVerification;

};

 struct reqStruct
{
    public string trasaID, tfaMethod, totpCode, user, userIP, workstation, serviceID, serviceKey;


}
namespace TrasaTfaPrompt
{

    public partial class Form1 : Form
    {

        public string userValFromCP { get; set; }
        public string ReturnValue1 { get; set; }

        public Form1()
        {
            InitializeComponent();
        }


        reqStruct reqVals;
        private CountDownTimer timer = new CountDownTimer();


        private void Label1_Click(object sender, EventArgs e)
        {

        }

        private void Form1_Load(object sender, EventArgs e)
        {
            //backgroundWorker1.RunWorkerAsync();
            this.timer.SetTime(0, 59);
            timer.Start();
            timer.TimeChanged += () => this.progressBar1.PerformStep();
            timer.CountDownFinished += () => this.exitWithCode1();
            

        }

        private void exitWithCode1()
        {
            this.ReturnValue1 = "failed";
            Application.Exit();
        }


        private void Label2_Click(object sender, EventArgs e)
        {

        }

        private void Label3_Click(object sender, EventArgs e)
        {

        }


        private void TextBox2_TextChanged(object sender, EventArgs e)
        {
            reqVals.totpCode = totpCode.Text;
        }

        private void SelectTfaMethod_SelectedIndexChanged(object sender, EventArgs e)
        {
            reqVals.tfaMethod = "U2F";
            string tfaMethod = selectTfaMethod.SelectedItem.ToString();
            int cmp = string.Compare(tfaMethod, "U2F");
            if (cmp == 0)
            {
                reqVals.tfaMethod = "U2F";
                totpCode.Visible = false;
            } else
            {
                reqVals.tfaMethod = "totp";
                totpCode.Visible = true;
            }
        }

        private void CloseBtn_Click(object sender, EventArgs e)
        {
            //Application.Exit();
              this.ReturnValue1 = "failed";
             //   shut.LogOff();
                Application.Exit();
        }

        private void TrasaEmail_TextChanged(object sender, EventArgs e)
        {
            reqVals.trasaID = trasaEmail.Text;
        }



    
        private void SubmitTfa_Click(object sender, EventArgs e)
        {

            trasaCPConfig configFile; //= JsonConvert.DeserializeObject<trasaCPConfig>(File.ReadAllText(@"C:\Program Files\TrasaWin\config.trasa"));
            string[] hosts = new string[] { };
            string[] users = new string[] { };
            try
            {
                configFile = JsonConvert.DeserializeObject<trasaCPConfig>(File.ReadAllText(@"C:\Program Files\TrasaWin\config.trasa"));
                reqVals.user = this.userValFromCP;
                reqVals.serviceID = configFile.serviceID;
                reqVals.serviceKey = configFile.serviceKey;
                configFile.trasaHost = configFile.trasaHost.Replace(" ", String.Empty);
                hosts = configFile.trasaHost.Split(',');
                configFile.offlineUsers = configFile.offlineUsers.Replace(" ", String.Empty);
                users = configFile.offlineUsers.Split(',');

            }
            catch (Exception je)
            {
                this.ReturnValue1 = "success";
                var now = DateTime.Now;
                string val = string.Format("{0} - Authentication was bypassed for user {1} because invalid config file found  - error: {2} \r\n", now, this.userValFromCP, je.ToString());
                CreateFolder.saveFile(val, true);
                Application.Exit();
                return;
            }





            if (String.IsNullOrEmpty(reqVals.trasaID))
            {
                MessageBox.Show("You need to submit your email or username");

                return;

            }

            statusText.Text = "Verifying request.....";
            this.Enabled = false;

            if (String.IsNullOrEmpty(reqVals.tfaMethod))
            {
                reqVals.tfaMethod = "U2F";
            }
            if (String.IsNullOrEmpty(reqVals.totpCode))
            {
                reqVals.totpCode = "";
            }
            if (String.IsNullOrEmpty(reqVals.userIP))
            {
                reqVals.userIP = "";
            }
         
            if (String.IsNullOrEmpty(reqVals.workstation))
            {
                reqVals.workstation = "";
            }

           
           
            string requestVar = JsonConvert.SerializeObject(reqVals);
            string respCode = "failed";
           
            responseVal resp;
            resp = checkResp.makeTfaRequest(hosts, requestVar, configFile.skipTLSVerification);

            // if the respCode is success we return success and exit the program.
            // if the respCode is failed, we return failed and exit the program.
            // if respCode is invalid, this means that trasaWin could not connect to any of trasacore hosts. This is offline case and we should check if the user is authorized in offline mode.
            // 
            if (resp.status.Equals("success"))
            {
                this.ReturnValue1 = "success";
                Application.Exit();
            } else if (resp.status.Equals("failed"))
            {
                statusText.BackColor = System.Drawing.Color.Maroon;
                statusText.ForeColor = System.Drawing.Color.White;
                statusText.Text = "Failed 2FA";
                this.Enabled = true;
                MessageBox.Show(resp.reason);
                // Prompt.ShowDialog(resp.reason);
                // this.ReturnValue1 = "failed";
                //shut.LogOff();
                //Application.Exit();
            } else
            {
                // we now check for offline users.
                respCode = checkOfflineUsers.check(users, this.userValFromCP);

                if (respCode.Equals("success"))
                {
                    this.ReturnValue1 = "success";
                    var now = DateTime.Now;
                    string val = string.Format("{0} - Allowing offline access to user {1} \r\n", now, this.userValFromCP);
                    Application.Exit();
                }
                else 
                {
                    this.ReturnValue1 = "invalid-offline";
                    MessageBox.Show("TrasWIN cannot contact trasacore (trasa server) and currently administrator has not authorized this user for offline access.");
                    var now = DateTime.Now;
                    string val = string.Format("{0} - Login blocked as trasa cannot contact trasacore and offline usage has not been authorized for user {1} \r\n", now, this.userValFromCP);
                    CreateFolder.saveFile(val, true);
                    // shut.LogOff();
                    Application.Exit();
                }
            }

         
        }

        private void ProgressBar1_Click(object sender, EventArgs e)
        {

        }

        private void StatusText_Click(object sender, EventArgs e)
        {

        }
    }


    public class shut
    {
        [DllImport("user32.dll", SetLastError = true)]
        static extern bool LockWorkStation();
       // static extern int ExitWindowsEx(uint uFlags, uint dwReason);
        
        // public static void LogOff() { }
        public static void LogOff()
        {
            LockWorkStation();
           // const ushort EWX_FORCE = 0x00000004;
            // getPrivileges();
           // ExitWindowsEx(EWX_FORCE, 0);
        }
    }


public class checkResp
    {
        public static responseVal makeTfaRequest(string[] hosts, string reqVal, bool skipVerify)
        {
            responseVal resp;
            resp.status = "invalid";
            resp.reason = "nth";
            resp.data = new string[6];
            resp.error = "nth";
            resp.intent = "";
            foreach (string host in hosts)
            {
                resp = Sendtfareq.Web.SendRequest(host, reqVal, skipVerify);
                if (resp.status.Equals("success"))
                {
                    //return resp;
                    break;
                } else if (resp.status.Equals("failed"))
                {
                    // return resp;
                    break;
                }
            }
          

            return resp;
        }
    }




    public class CountDownTimer : IDisposable
    {
        public Action TimeChanged;
        public Action CountDownFinished;

        public bool IsRunnign => timer.Enabled;

        public int StepMs
        {
            get => timer.Interval;
            set => timer.Interval = value;
        }

        private Timer timer = new Timer();

        private DateTime _maxTime = new DateTime(1, 1, 1, 0, 30, 0);
        private DateTime _minTime = new DateTime(1, 1, 1, 0, 0, 0);

        public DateTime TimeLeft { get; private set; }
        private long TimeLeftMs => TimeLeft.Ticks / TimeSpan.TicksPerMillisecond;

        public string TimeLeftStr => TimeLeft.ToString("mm:ss");

        public string TimeLeftMsStr => TimeLeft.ToString("mm:ss.fff");

        private void TimerTick(object sender, EventArgs e)
        {
            if (TimeLeftMs > timer.Interval)
            {
                TimeLeft = TimeLeft.AddMilliseconds(-timer.Interval);
                TimeChanged?.Invoke();
            }
            else
            {
                Stop();
                TimeLeft = _minTime;

                TimeChanged?.Invoke();
                CountDownFinished?.Invoke();
            }
        }

        public CountDownTimer(int min, int sec)
        {
            SetTime(min, sec);
            Init();
        }

        public CountDownTimer(DateTime dt)
        {
            SetTime(dt);
            Init();
        }

        public CountDownTimer()
        {
            Init();
        }

        private void Init()
        {
            TimeLeft = _maxTime;

            StepMs = 1000;
            timer.Tick += new EventHandler(TimerTick);
        }

        public void SetTime(DateTime dt)
        {
            TimeLeft = _maxTime = dt;
            TimeChanged?.Invoke();
        }

        public void SetTime(int min, int sec = 0) => SetTime(new DateTime(1, 1, 1, 0, min, sec));

        public void Start() => timer.Start();

        public void Pause() => timer.Stop();

        public void Stop()
        {
            Pause();
            Reset();
        }

        public void Reset()
        {
            TimeLeft = _maxTime;
        }

        public void Restart()
        {
            Reset();
            Start();
        }

        public void Dispose() => timer.Dispose();
    }



    public class checkOfflineUsers
    {
        public static string check(string[] users, string user)
        {
            string respCode = "failed";
            foreach (string u in users)
            {

                if (u.Equals(user))
                {
                    respCode = "success";
                    break;
                }
                else
                {
                    respCode = "failed";
                }
            }
            return respCode;
        }
    }

   public struct responseVal
    {
        public string status, reason, intent, error ;
        public string[] data;
    };

    namespace Sendtfareq
    {
  
        public class Web
        {
          
            public static responseVal SendRequest(string reqMeta, string requestData, bool skipVerify)
            {
                responseVal responseData;
                string endpoint = reqMeta + "/auth/agent/win";

                responseData.status = "failed";
                responseData.intent = "";
                    responseData.data = new string[6];
                    responseData.error = "";
                    responseData.reason = "";

                   

                try
                {
                    //  Prompt.ShowDialog(reqMeta);
                    // Create a request using a URL that can receive a post.   
                  //   WebRequest request = WebRequest.Create(endpoint);

                    HttpWebRequest request = HttpWebRequest.CreateHttp(endpoint);

                    if (skipVerify) {
                        request.ServerCertificateValidationCallback += (sender, certificate, chain, sslPolicyErrors) => true;
                    }

                    // Set the Method property of the request to POST.  
                    request.Method = "POST";
                    // Create POST data and convert it to a byte array.  
                    string postData = requestData;
                    byte[] byteArray = Encoding.UTF8.GetBytes(postData);
                    // Set the ContentType property of the WebRequest.  
                    request.ContentType = "application/json";
                    // Set the ContentLength property of the WebRequest.  
                    request.ContentLength = byteArray.Length;
                    // Get the request stream.  
                    Stream dataStream = request.GetRequestStream();
                    // Write the data to the request stream.  
                    dataStream.Write(byteArray, 0, byteArray.Length);
                    // Close the Stream object.  
                    dataStream.Close();

                    // Get the response.  
                    WebResponse response = request.GetResponse();
                    // Display the status.  
                    // Console.WriteLine(((HttpWebResponse)response).StatusDescription);
                    // Get the stream containing content returned by the server.  
                    dataStream = response.GetResponseStream();
                    // Open the stream using a StreamReader for easy access.  
                    StreamReader reader = new StreamReader(dataStream);
                    // Read the content.  
                    string responseFromServer = reader.ReadToEnd();
                    // Display the content.  
                    // Console.WriteLine(responseFromServer);

                    // responseVal responseData ;
                    
                    responseData = JsonConvert.DeserializeObject<responseVal>(responseFromServer);
                    reader.Close();
                    dataStream.Close();
                    response.Close();
                    return responseData;


                    //int compare = responseStruct.status.CompareTo("success");
                    //if (compare == 0)
                    //{
                    //    reader.Close();
                    //    dataStream.Close();
                    //    response.Close();
                    //    return "success";
                    //}

                    //else
                    //{
                    //    // Clean up the streams.  
                    //    reader.Close();
                    //    dataStream.Close();
                    //    response.Close();
                    //    return "failed";
                    //}
                }
                catch (WebException webExcp)
                {
                    var now = DateTime.Now;
                    string val = string.Format("{0} - Failed to contact server: {1} \r\n", now, reqMeta);
                    CreateFolder.saveFile(val, true);
                    responseData.status = "invalid";
                    
                    WebExceptionStatus status = webExcp.Status;
                    if (status == WebExceptionStatus.ProtocolError)
                    {
                       
                        return responseData;
                    }
                    return  responseData;
                }
                catch (Exception e)
                {
                    responseData.status = e.ToString();
                    return responseData;
                }



            }
        }
    }

    public class CreateFolder
    {
        public static void saveFile(string configValues, bool update)
        {

            // You can write out the path name directly instead of using the Combine
            // method. Combine just makes the process easier.
            string pathString = @"C:\Program Files\TrasaWin";


            // Create the subfolder. You can verify in File Explorer that you have this
            // structure in the C: drive.
            //    Local Disk (C:)
            //        Top-Level Folder
            //            SubFolder
            System.IO.Directory.CreateDirectory(pathString);

            // Create a file name for the file you want to create. 
            string fileName = "TrasaLogs.txt";

            // This example uses a random string for the name, but you also can specify
            // a particular name.
            //string fileName = "MyNewFile.txt";

            // Use Combine again to add the file name to the path.
            pathString = System.IO.Path.Combine(pathString, fileName);

            // Verify the path that you have constructed.
           // Console.WriteLine("Path to my file: {0}\n", pathString);

            // Check that the file doesn't already exist. If it doesn't exist, create
            // the file and write integers 0 - 99 to it.
            // DANGER: System.IO.File.Create will overwrite the file if it already exists.
            // This could happen even with random file names, although it is unlikely.


            using (System.IO.StreamWriter fs = System.IO.File.AppendText(pathString))
            {

                AddText(fs, configValues, update);


            }

        }

        private static void AddText(StreamWriter fs, string value, bool update)
        {
            update = false;
            Log(value, fs);

        }

        public static void Log(string logMessage, TextWriter w)
        {
            w.WriteLine($"{logMessage}\r\n");
          
        }
    }
    public static class Prompt
    {
        public static int ShowDialog(string text)
        {
            Form prompt = new Form();
            prompt.StartPosition = System.Windows.Forms.FormStartPosition.CenterScreen;
            prompt.Width = 400;
            prompt.Height = 200;
            prompt.ShowIcon = false;

            Label textLabel = new Label(); // { Left = 50, Top = 20 };

            textLabel.AutoSize = true;
            textLabel.Font = new System.Drawing.Font("Microsoft Sans Serif", 9F, System.Drawing.FontStyle.Regular, System.Drawing.GraphicsUnit.Point, ((byte)(0)));
            textLabel.Location = new System.Drawing.Point(24, 50);
            textLabel.Size = new System.Drawing.Size(152, 20);
            textLabel.TabIndex = 0;
            textLabel.Text = text;


            // text tests = new text() 
            // TextBox textBox = new TextBox() { Left = 50, Top = 50, Width = 200, Text = text };
            //NumericUpDown inputBox = new NumericUpDown() { Left = 50, Top = 50, Width = 400 };
            Button confirmation = new Button() { Text = "Ok", Left = 150, Width = 100, Top = 100 };
            confirmation.Click += (sender, e) => { prompt.Close(); };
            prompt.Controls.Add(confirmation);
            prompt.Controls.Add(textLabel);
            // prompt.Controls.Add(textBox);
            // prompt.Controls.Add(inputBox);
            prompt.ShowDialog();
            return 0;
        }
    }
}



