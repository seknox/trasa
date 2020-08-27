
namespace TrasaConfig
{
    partial class Form1
    {
        /// <summary>
        /// Required designer variable.
        /// </summary>
        private System.ComponentModel.IContainer components = null;

        /// <summary>
        /// Clean up any resources being used.
        /// </summary>
        /// <param name="disposing">true if managed resources should be disposed; otherwise, false.</param>
        protected override void Dispose(bool disposing)
        {
            if (disposing && (components != null))
            {
                components.Dispose();
            }
            base.Dispose(disposing);
        }

        #region Windows Form Designer generated code

        /// <summary>
        /// Required method for Designer support - do not modify
        /// the contents of this method with the code editor.
        /// </summary>
        /// 
     
        private void InitializeComponent()
        {
            System.ComponentModel.ComponentResourceManager resources = new System.ComponentModel.ComponentResourceManager(typeof(Form1));
            this.save_button = new System.Windows.Forms.Button();
            this.label8 = new System.Windows.Forms.Label();
            this.appConfig_tab = new System.Windows.Forms.TabPage();
            this.skipTLSVerification = new System.Windows.Forms.CheckBox();
            this.OfflineUserNames = new System.Windows.Forms.TextBox();
            this.OfflineUsers = new System.Windows.Forms.Label();
            this.serviceKey_tb = new System.Windows.Forms.TextBox();
            this.trasaHost_tb = new System.Windows.Forms.TextBox();
            this.serviceID_tb = new System.Windows.Forms.TextBox();
            this.serviceKey_lbl = new System.Windows.Forms.Label();
            this.trasaHost_lbl = new System.Windows.Forms.Label();
            this.serviceID_lbl = new System.Windows.Forms.Label();
            this.settingTab = new System.Windows.Forms.TabControl();
            this.appConfig_tab.SuspendLayout();
            this.settingTab.SuspendLayout();
            this.SuspendLayout();
            // 
            // save_button
            // 
            this.save_button.BackColor = System.Drawing.Color.DodgerBlue;
            this.save_button.BackgroundImageLayout = System.Windows.Forms.ImageLayout.Center;
            this.save_button.Cursor = System.Windows.Forms.Cursors.Hand;
            this.save_button.FlatStyle = System.Windows.Forms.FlatStyle.Popup;
            this.save_button.Font = new System.Drawing.Font("Microsoft Sans Serif", 9.75F, System.Drawing.FontStyle.Bold, System.Drawing.GraphicsUnit.Point, ((byte)(0)));
            this.save_button.ForeColor = System.Drawing.SystemColors.ButtonHighlight;
            this.save_button.Location = new System.Drawing.Point(232, 256);
            this.save_button.Name = "save_button";
            this.save_button.Size = new System.Drawing.Size(159, 31);
            this.save_button.TabIndex = 1;
            this.save_button.Text = "Save Configuration";
            this.save_button.UseVisualStyleBackColor = false;
            this.save_button.Click += new System.EventHandler(this.save_button_Click);
            // 
            // label8
            // 
            this.label8.AutoSize = true;
            this.label8.Font = new System.Drawing.Font("Microsoft Sans Serif", 8.25F, System.Drawing.FontStyle.Regular, System.Drawing.GraphicsUnit.Point, ((byte)(0)));
            this.label8.Location = new System.Drawing.Point(12, 462);
            this.label8.Name = "label8";
            this.label8.Size = new System.Drawing.Size(264, 13);
            this.label8.TabIndex = 5;
            this.label8.Text = "TrasaWIN, proudly a product of Seknox  Cybersecurity";
            // 
            // appConfig_tab
            // 
            this.appConfig_tab.Controls.Add(this.skipTLSVerification);
            this.appConfig_tab.Controls.Add(this.save_button);
            this.appConfig_tab.Controls.Add(this.OfflineUserNames);
            this.appConfig_tab.Controls.Add(this.OfflineUsers);
            this.appConfig_tab.Controls.Add(this.serviceKey_tb);
            this.appConfig_tab.Controls.Add(this.trasaHost_tb);
            this.appConfig_tab.Controls.Add(this.serviceID_tb);
            this.appConfig_tab.Controls.Add(this.serviceKey_lbl);
            this.appConfig_tab.Controls.Add(this.trasaHost_lbl);
            this.appConfig_tab.Controls.Add(this.serviceID_lbl);
            this.appConfig_tab.Location = new System.Drawing.Point(4, 25);
            this.appConfig_tab.Name = "appConfig_tab";
            this.appConfig_tab.Padding = new System.Windows.Forms.Padding(3);
            this.appConfig_tab.Size = new System.Drawing.Size(579, 394);
            this.appConfig_tab.TabIndex = 0;
            this.appConfig_tab.Text = "Configure TRASA 2FA agent";
            this.appConfig_tab.UseVisualStyleBackColor = true;
            // 
            // skipTLSVerification
            // 
            this.skipTLSVerification.AutoSize = true;
            this.skipTLSVerification.Location = new System.Drawing.Point(232, 215);
            this.skipTLSVerification.Name = "skipTLSVerification";
            this.skipTLSVerification.Size = new System.Drawing.Size(151, 20);
            this.skipTLSVerification.TabIndex = 8;
            this.skipTLSVerification.Text = "Skip TLS Verification";
            this.skipTLSVerification.UseVisualStyleBackColor = true;
            this.skipTLSVerification.CheckedChanged += new System.EventHandler(this.skipTLSVerification_CheckedChanged);
            // 
            // OfflineUserNames
            // 
            this.OfflineUserNames.Font = new System.Drawing.Font("Microsoft Sans Serif", 8.25F, System.Drawing.FontStyle.Regular, System.Drawing.GraphicsUnit.Point, ((byte)(0)));
            this.OfflineUserNames.Location = new System.Drawing.Point(182, 168);
            this.OfflineUserNames.Name = "OfflineUserNames";
            this.OfflineUserNames.Size = new System.Drawing.Size(353, 20);
            this.OfflineUserNames.TabIndex = 7;
            this.OfflineUserNames.TextChanged += new System.EventHandler(this.OfflineUserNames_TextChanged);
            // 
            // OfflineUsers
            // 
            this.OfflineUsers.AutoSize = true;
            this.OfflineUsers.Font = new System.Drawing.Font("Microsoft Sans Serif", 12F, System.Drawing.FontStyle.Regular, System.Drawing.GraphicsUnit.Point, ((byte)(0)));
            this.OfflineUsers.Location = new System.Drawing.Point(26, 166);
            this.OfflineUsers.Name = "OfflineUsers";
            this.OfflineUsers.Size = new System.Drawing.Size(109, 20);
            this.OfflineUsers.TabIndex = 6;
            this.OfflineUsers.Text = "Offline Users :";
            // 
            // serviceKey_tb
            // 
            this.serviceKey_tb.Font = new System.Drawing.Font("Microsoft Sans Serif", 8.25F, System.Drawing.FontStyle.Regular, System.Drawing.GraphicsUnit.Point, ((byte)(0)));
            this.serviceKey_tb.Location = new System.Drawing.Point(182, 91);
            this.serviceKey_tb.Name = "serviceKey_tb";
            this.serviceKey_tb.Size = new System.Drawing.Size(354, 20);
            this.serviceKey_tb.TabIndex = 5;
            this.serviceKey_tb.TextChanged += new System.EventHandler(this.serviceKey_tb_TextChanged);
            // 
            // trasaHost_tb
            // 
            this.trasaHost_tb.Font = new System.Drawing.Font("Microsoft Sans Serif", 8.25F, System.Drawing.FontStyle.Regular, System.Drawing.GraphicsUnit.Point, ((byte)(0)));
            this.trasaHost_tb.Location = new System.Drawing.Point(182, 133);
            this.trasaHost_tb.Name = "trasaHost_tb";
            this.trasaHost_tb.Size = new System.Drawing.Size(354, 20);
            this.trasaHost_tb.TabIndex = 4;
            this.trasaHost_tb.TextChanged += new System.EventHandler(this.trasaHost_tb_TextChanged);
            // 
            // serviceID_tb
            // 
            this.serviceID_tb.Font = new System.Drawing.Font("Microsoft Sans Serif", 8.25F, System.Drawing.FontStyle.Regular, System.Drawing.GraphicsUnit.Point, ((byte)(0)));
            this.serviceID_tb.Location = new System.Drawing.Point(182, 50);
            this.serviceID_tb.Name = "serviceID_tb";
            this.serviceID_tb.Size = new System.Drawing.Size(354, 20);
            this.serviceID_tb.TabIndex = 1;
            this.serviceID_tb.TextChanged += new System.EventHandler(this.serviceID_tb_TextChanged);
            // 
            // serviceKey_lbl
            // 
            this.serviceKey_lbl.AutoSize = true;
            this.serviceKey_lbl.Font = new System.Drawing.Font("Microsoft Sans Serif", 12F, System.Drawing.FontStyle.Regular, System.Drawing.GraphicsUnit.Point, ((byte)(0)));
            this.serviceKey_lbl.Location = new System.Drawing.Point(37, 90);
            this.serviceKey_lbl.Name = "serviceKey_lbl";
            this.serviceKey_lbl.Size = new System.Drawing.Size(95, 20);
            this.serviceKey_lbl.TabIndex = 3;
            this.serviceKey_lbl.Text = "ServiceKey :";
            // 
            // trasaHost_lbl
            // 
            this.trasaHost_lbl.AutoSize = true;
            this.trasaHost_lbl.Font = new System.Drawing.Font("Microsoft Sans Serif", 12F, System.Drawing.FontStyle.Regular, System.Drawing.GraphicsUnit.Point, ((byte)(0)));
            this.trasaHost_lbl.Location = new System.Drawing.Point(5, 131);
            this.trasaHost_lbl.Name = "trasaHost_lbl";
            this.trasaHost_lbl.Size = new System.Drawing.Size(132, 20);
            this.trasaHost_lbl.TabIndex = 2;
            this.trasaHost_lbl.Text = "TRASA Host/IP : ";
            // 
            // serviceID_lbl
            // 
            this.serviceID_lbl.AutoSize = true;
            this.serviceID_lbl.Font = new System.Drawing.Font("Microsoft Sans Serif", 12F, System.Drawing.FontStyle.Regular, System.Drawing.GraphicsUnit.Point, ((byte)(0)));
            this.serviceID_lbl.Location = new System.Drawing.Point(43, 50);
            this.serviceID_lbl.Name = "serviceID_lbl";
            this.serviceID_lbl.Size = new System.Drawing.Size(90, 20);
            this.serviceID_lbl.TabIndex = 0;
            this.serviceID_lbl.Text = "ServiceID : ";
            this.serviceID_lbl.Click += new System.EventHandler(this.serviceID_lbl_Click);
            // 
            // settingTab
            // 
            this.settingTab.AccessibleName = "Settings";
            this.settingTab.Controls.Add(this.appConfig_tab);
            this.settingTab.Font = new System.Drawing.Font("Microsoft Sans Serif", 9.75F, System.Drawing.FontStyle.Regular, System.Drawing.GraphicsUnit.Point, ((byte)(0)));
            this.settingTab.Location = new System.Drawing.Point(10, 20);
            this.settingTab.Name = "settingTab";
            this.settingTab.SelectedIndex = 0;
            this.settingTab.Size = new System.Drawing.Size(587, 423);
            this.settingTab.TabIndex = 0;
            // 
            // Form1
            // 
            this.AutoScaleDimensions = new System.Drawing.SizeF(6F, 13F);
            this.AutoScaleMode = System.Windows.Forms.AutoScaleMode.Font;
            this.BackColor = System.Drawing.SystemColors.ControlLightLight;
            this.ClientSize = new System.Drawing.Size(604, 505);
            this.Controls.Add(this.label8);
            this.Controls.Add(this.settingTab);
            this.Icon = ((System.Drawing.Icon)(resources.GetObject("$this.Icon")));
            this.Name = "Form1";
            this.StartPosition = System.Windows.Forms.FormStartPosition.CenterScreen;
            this.Text = "TrasaWIN";
            this.TopMost = true;
            this.Load += new System.EventHandler(this.Form1_Load);
            this.appConfig_tab.ResumeLayout(false);
            this.appConfig_tab.PerformLayout();
            this.settingTab.ResumeLayout(false);
            this.ResumeLayout(false);
            this.PerformLayout();

        }

        #endregion
        private System.Windows.Forms.Button save_button;
        private System.Windows.Forms.Label label8;
        private System.Windows.Forms.TabPage appConfig_tab;
        private System.Windows.Forms.TextBox serviceKey_tb;
        private System.Windows.Forms.TextBox trasaHost_tb;
        private System.Windows.Forms.TextBox serviceID_tb;
        private System.Windows.Forms.Label serviceKey_lbl;
        private System.Windows.Forms.Label trasaHost_lbl;
        private System.Windows.Forms.Label serviceID_lbl;
        private System.Windows.Forms.TabControl settingTab;
        private System.Windows.Forms.TextBox OfflineUserNames;
        private System.Windows.Forms.Label OfflineUsers;
        private System.Windows.Forms.CheckBox skipTLSVerification;
    }
}

