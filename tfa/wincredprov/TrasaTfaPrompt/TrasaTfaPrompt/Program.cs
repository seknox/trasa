using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using System.Windows.Forms;

using RGiesecke.DllExport;
using System.Runtime.InteropServices;

namespace TrasaTfaPrompt
{
    
    static class Program
    {
     
        /// <summary>
        /// The main entry point for the application.
        /// </summary>

        [DllExport(ExportName = "TrasaTfaPrompt", CallingConvention = CallingConvention.StdCall)]
        [STAThread]
        [return: MarshalAs(UnmanagedType.LPWStr)]
        public static string ShowForm([MarshalAs(UnmanagedType.LPWStr)]string userName)
        {
              
            Form1 form = new Form1();
            form.userValFromCP = userName;
           form.ShowDialog();
           return form.ReturnValue1;
         
          //  return "testval";
        }

       

        //[STAThread]
        //static void Main()
        //{

        //    Application.EnableVisualStyles();
        //    Application.SetCompatibleTextRenderingDefault(false);
        //    Application.Run(new Form1());
        //}
    }
}


