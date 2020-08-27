using System;
using System.Text;
using System.Windows.Forms;
using Newtonsoft.Json;
using System.IO;
using System.Net;

//using System.Net.Http;
//using System.Text;


struct trasaCPConfig
{
    public string serviceID, serviceKey, trasaHost, offlineUsers; public bool skipTLSVerification;
};


namespace TrasaConfig
{
    public partial class Form1 : Form
    {
        public Form1()
        {
            InitializeComponent();
            trasaCPConfig storedConfig;
            try
            {
                storedConfig = JsonConvert.DeserializeObject<trasaCPConfig>(File.ReadAllText(@"C:\Program Files\TrasaWin\config.trasa"));
                if (storedConfig.serviceID != null)
                {
                    serviceID_tb.Text = storedConfig.serviceID;
                }
                if (storedConfig.serviceKey != null)
                {
                    serviceKey_tb.Text = storedConfig.serviceKey;
                }
                if (storedConfig.trasaHost != null)
                {
                    trasaHost_tb.Text = storedConfig.trasaHost;
                }
                if (storedConfig.offlineUsers != null)
                {
                    OfflineUserNames.Text = storedConfig.offlineUsers;
                }
                if (storedConfig.skipTLSVerification)
                {
                    skipTLSVerification.Checked = storedConfig.skipTLSVerification;


                }
            }
            catch(Exception je)
            {
                Console.WriteLine(je);// return;
            }

           

        }


        trasaCPConfig configvalues;

        private void trasaHost_tb_TextChanged(object sender, EventArgs e)
        {
            configvalues.trasaHost = trasaHost_tb.Text;
        }

        private void serviceKey_tb_TextChanged(object sender, EventArgs e)
        {
            configvalues.serviceKey = serviceKey_tb.Text;
        }

        private void serviceID_tb_TextChanged(object sender, EventArgs e)
        {
            configvalues.serviceID = serviceID_tb.Text;
        }

              private void OfflineUserNames_TextChanged(object sender, EventArgs e)
        {
            configvalues.offlineUsers = OfflineUserNames.Text;
        }

        private void skipTLSVerification_CheckedChanged(object sender, EventArgs e)
        {
           
            configvalues.skipTLSVerification = skipTLSVerification.Checked;
            // MessageBox.Show("clicked: ");
        }

        private void save_button_Click(object sender, EventArgs e)
        {
            if (configvalues.offlineUsers == null )
            {
                configvalues.offlineUsers = "";
            }
            string requestVar = JsonConvert.SerializeObject(configvalues);
            string[] url = configvalues.trasaHost.Split(',');
            string host = string.Format("{0}/auth/agent/checkconfig", url);
            string resp = Sendtfareq.Web.SendRequest(host, requestVar, configvalues.skipTLSVerification);
            if (resp.Equals("success"))
            {
                string output = JsonConvert.SerializeObject(configvalues);
                CreateFolder.saveFile(output);
                MessageBox.Show("Successfully verified and saved config data.");
               
               // MessageBox.Show("Successfully stored config data.");
                //  Prompt.ShowDialog("Successfully verified. You can save your settings now.");
            }
            else if (resp.Equals("invalid"))
            {
                MessageBox.Show("Failed. TRASA host value is invalid. DO not exit the application until config values are successfully verified.");
                //  Prompt.ShowDialog("Failed. TRASA host value is invalid");
            }
            else
            {
                MessageBox.Show("Failed. Please check your configurations again. DO not exit the application until config values are successfully verified");
                //  Prompt.ShowDialog("Failed. Please check your configurations again");
            }

          
          //  Prompt.ShowDialog("Successfully stored config data.");
        }


        private void Form1_Load(object sender, EventArgs e)
        {

        }

        private void serviceID_lbl_Click(object sender, EventArgs e)
        {

        }

        private void TabPage1_Click(object sender, EventArgs e)
        {

        }

        private void Label2_Click(object sender, EventArgs e)
        {

        }

  
    }


    public class CreateFolder
    {
       public static void saveFile(string configValues)
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
            string fileName = "config.trasa";

            // This example uses a random string for the name, but you also can specify
            // a particular name.
            //string fileName = "MyNewFile.txt";

            // Use Combine again to add the file name to the path.
            pathString = System.IO.Path.Combine(pathString, fileName);

            // Verify the path that you have constructed.
            Console.WriteLine("Path to my file: {0}\n", pathString);

            // Check that the file doesn't already exist. If it doesn't exist, create
            // the file and write integers 0 - 99 to it.
            // DANGER: System.IO.File.Create will overwrite the file if it already exists.
            // This could happen even with random file names, although it is unlikely.

            
                using (System.IO.FileStream fs = System.IO.File.Create(pathString))
                {

                AddText(fs, configValues); 


                }
           
        }

        private static void AddText(FileStream fs, string value)
        {
            byte[] info = new UTF8Encoding(true).GetBytes(value);
            fs.Write(info, 0, info.Length);
        }
    }

    struct responseVal
    {
        public string status;
    };


    namespace Sendtfareq
    {
        public class Web
        {
            public static string SendRequest(string reqMeta, string requestData, bool skipVerify)
            {
                
                //string endpoint = "https://" + host + "/api/v1/remote/auth/checkconfig";
           
                try
                {
                    //  Prompt.ShowDialog(reqMeta);
                    // Create a request using a URL that can receive a post.   
                    // WebRequest request = WebRequest.Create(reqMeta);
                    HttpWebRequest request = HttpWebRequest.CreateHttp(reqMeta);

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
                    responseVal responseStruct;
                    responseVal responseData = JsonConvert.DeserializeObject<responseVal>(responseFromServer);

                    responseStruct.status = responseData.status;

                    int compare = responseStruct.status.CompareTo("success");
                    if (compare == 0)
                    {
                        Console.WriteLine("Successfully verified :)", "SUCCESS");
                        reader.Close();
                        dataStream.Close();
                        response.Close();
                        return "success";
                    }

                    else
                    {

                        Console.WriteLine("Check configuration values again", "FAILED");
                        // Clean up the streams.  
                        reader.Close();
                        dataStream.Close();
                        response.Close();
                        return "failed";
                    }
                }
                catch (WebException webExcp)
                {
                   
                    WebExceptionStatus status = webExcp.Status;
                    if (status == WebExceptionStatus.ProtocolError)
                    {
                        return "invalid";
                    }
                    return "invalid";
                }
                catch (Exception e)
                {
                    return e.ToString();
                }

               

            }
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
