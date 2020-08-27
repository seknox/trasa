namespace TrasaTfaPrompt
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
        private void InitializeComponent()
        {
            System.ComponentModel.ComponentResourceManager resources = new System.ComponentModel.ComponentResourceManager(typeof(Form1));
            this.label1 = new System.Windows.Forms.Label();
            this.label2 = new System.Windows.Forms.Label();
            this.totpCode = new System.Windows.Forms.TextBox();
            this.submitTfa = new System.Windows.Forms.Button();
            this.pictureBox1 = new System.Windows.Forms.PictureBox();
            this.groupBox1 = new System.Windows.Forms.GroupBox();
            this.trasaEmail = new System.Windows.Forms.TextBox();
            this.selectTfaMethod = new System.Windows.Forms.ComboBox();
            this.label3 = new System.Windows.Forms.Label();
            this.closeBtn = new System.Windows.Forms.Button();
            this.label4 = new System.Windows.Forms.Label();
            this.progressBar1 = new System.Windows.Forms.ProgressBar();
            this.status = new System.Windows.Forms.Label();
            this.statusText = new System.Windows.Forms.Label();
            ((System.ComponentModel.ISupportInitialize)(this.pictureBox1)).BeginInit();
            this.groupBox1.SuspendLayout();
            this.SuspendLayout();
            // 
            // label1
            // 
            this.label1.AutoSize = true;
            this.label1.Font = new System.Drawing.Font("Microsoft Sans Serif", 12F, System.Drawing.FontStyle.Bold, System.Drawing.GraphicsUnit.Point, ((byte)(0)));
            this.label1.Location = new System.Drawing.Point(495, 87);
            this.label1.Name = "label1";
            this.label1.Size = new System.Drawing.Size(377, 20);
            this.label1.TabIndex = 0;
            this.label1.Text = "Enter your TRASA username or email address";
            this.label1.Click += new System.EventHandler(this.Label1_Click);
            // 
            // label2
            // 
            this.label2.AutoSize = true;
            this.label2.Font = new System.Drawing.Font("Microsoft Sans Serif", 12F, System.Drawing.FontStyle.Bold, System.Drawing.GraphicsUnit.Point, ((byte)(0)));
            this.label2.Location = new System.Drawing.Point(497, 154);
            this.label2.Name = "label2";
            this.label2.Size = new System.Drawing.Size(173, 20);
            this.label2.TabIndex = 2;
            this.label2.Text = "Choose 2FA Method";
            this.label2.Click += new System.EventHandler(this.Label2_Click);
            // 
            // totpCode
            // 
            this.totpCode.Location = new System.Drawing.Point(33, 163);
            this.totpCode.MaxLength = 6;
            this.totpCode.Name = "totpCode";
            this.totpCode.Size = new System.Drawing.Size(152, 20);
            this.totpCode.TabIndex = 3;
            this.totpCode.Visible = false;
            this.totpCode.TextChanged += new System.EventHandler(this.TextBox2_TextChanged);
            // 
            // submitTfa
            // 
            this.submitTfa.BackColor = System.Drawing.Color.MidnightBlue;
            this.submitTfa.Font = new System.Drawing.Font("Microsoft Sans Serif", 15.75F, System.Drawing.FontStyle.Bold, System.Drawing.GraphicsUnit.Point, ((byte)(0)));
            this.submitTfa.ForeColor = System.Drawing.SystemColors.ButtonHighlight;
            this.submitTfa.Location = new System.Drawing.Point(501, 244);
            this.submitTfa.Name = "submitTfa";
            this.submitTfa.Size = new System.Drawing.Size(160, 34);
            this.submitTfa.TabIndex = 4;
            this.submitTfa.Text = "Submit";
            this.submitTfa.UseVisualStyleBackColor = false;
            this.submitTfa.Click += new System.EventHandler(this.SubmitTfa_Click);
            // 
            // pictureBox1
            // 
            this.pictureBox1.BackgroundImageLayout = System.Windows.Forms.ImageLayout.Center;
            this.pictureBox1.Image = ((System.Drawing.Image)(resources.GetObject("pictureBox1.Image")));
            this.pictureBox1.Location = new System.Drawing.Point(22, 66);
            this.pictureBox1.Name = "pictureBox1";
            this.pictureBox1.Size = new System.Drawing.Size(383, 66);
            this.pictureBox1.TabIndex = 5;
            this.pictureBox1.TabStop = false;
            // 
            // groupBox1
            // 
            this.groupBox1.Controls.Add(this.trasaEmail);
            this.groupBox1.Controls.Add(this.selectTfaMethod);
            this.groupBox1.Controls.Add(this.totpCode);
            this.groupBox1.FlatStyle = System.Windows.Forms.FlatStyle.Popup;
            this.groupBox1.Location = new System.Drawing.Point(468, 41);
            this.groupBox1.Name = "groupBox1";
            this.groupBox1.Size = new System.Drawing.Size(409, 270);
            this.groupBox1.TabIndex = 7;
            this.groupBox1.TabStop = false;
            // 
            // trasaEmail
            // 
            this.trasaEmail.Location = new System.Drawing.Point(31, 70);
            this.trasaEmail.Name = "trasaEmail";
            this.trasaEmail.Size = new System.Drawing.Size(305, 20);
            this.trasaEmail.TabIndex = 5;
            this.trasaEmail.TextChanged += new System.EventHandler(this.TrasaEmail_TextChanged);
            // 
            // selectTfaMethod
            // 
            this.selectTfaMethod.FormattingEnabled = true;
            this.selectTfaMethod.Items.AddRange(new object[] {
            "U2F",
            "TOTP"});
            this.selectTfaMethod.Location = new System.Drawing.Point(33, 137);
            this.selectTfaMethod.Name = "selectTfaMethod";
            this.selectTfaMethod.Size = new System.Drawing.Size(152, 21);
            this.selectTfaMethod.TabIndex = 4;
            this.selectTfaMethod.SelectedIndexChanged += new System.EventHandler(this.SelectTfaMethod_SelectedIndexChanged);
            // 
            // label3
            // 
            this.label3.Font = new System.Drawing.Font("Tahoma", 9.75F, System.Drawing.FontStyle.Bold, System.Drawing.GraphicsUnit.Point, ((byte)(0)));
            this.label3.Location = new System.Drawing.Point(19, 158);
            this.label3.Name = "label3";
            this.label3.Size = new System.Drawing.Size(386, 89);
            this.label3.TabIndex = 8;
            this.label3.Text = "This system is protected by TRASA. If you don\'t have access, contact your adminis" +
    "trator. ";
            this.label3.Click += new System.EventHandler(this.Label3_Click);
            // 
            // closeBtn
            // 
            this.closeBtn.Location = new System.Drawing.Point(802, 415);
            this.closeBtn.Name = "closeBtn";
            this.closeBtn.Size = new System.Drawing.Size(75, 23);
            this.closeBtn.TabIndex = 9;
            this.closeBtn.Text = "CLOSE";
            this.closeBtn.UseVisualStyleBackColor = true;
            this.closeBtn.Click += new System.EventHandler(this.CloseBtn_Click);
            // 
            // label4
            // 
            this.label4.AutoSize = true;
            this.label4.Location = new System.Drawing.Point(341, 394);
            this.label4.Name = "label4";
            this.label4.Size = new System.Drawing.Size(197, 13);
            this.label4.TabIndex = 11;
            this.label4.Text = "This prompt will be closed in 60 seconds";
            // 
            // progressBar1
            // 
            this.progressBar1.ForeColor = System.Drawing.Color.Lime;
            this.progressBar1.Location = new System.Drawing.Point(22, 378);
            this.progressBar1.Maximum = 60;
            this.progressBar1.Name = "progressBar1";
            this.progressBar1.Size = new System.Drawing.Size(855, 13);
            this.progressBar1.Step = 1;
            this.progressBar1.Style = System.Windows.Forms.ProgressBarStyle.Continuous;
            this.progressBar1.TabIndex = 12;
            // 
            // status
            // 
            this.status.AutoSize = true;
            this.status.Font = new System.Drawing.Font("Palatino Linotype", 12F, System.Drawing.FontStyle.Bold, System.Drawing.GraphicsUnit.Point, ((byte)(0)));
            this.status.Location = new System.Drawing.Point(41, 306);
            this.status.Name = "status";
            this.status.Size = new System.Drawing.Size(86, 22);
            this.status.TabIndex = 13;
            this.status.Text = "STATUS : ";
            // 
            // statusText
            // 
            this.statusText.AutoSize = true;
            this.statusText.Font = new System.Drawing.Font("Trebuchet MS", 9.75F, System.Drawing.FontStyle.Bold, System.Drawing.GraphicsUnit.Point, ((byte)(0)));
            this.statusText.Location = new System.Drawing.Point(133, 309);
            this.statusText.Name = "statusText";
            this.statusText.Size = new System.Drawing.Size(90, 18);
            this.statusText.TabIndex = 14;
            this.statusText.Text = "TFA Required";
            this.statusText.Click += new System.EventHandler(this.StatusText_Click);
            // 
            // Form1
            // 
            this.AutoScaleDimensions = new System.Drawing.SizeF(6F, 13F);
            this.AutoScaleMode = System.Windows.Forms.AutoScaleMode.Font;
            this.BackColor = System.Drawing.SystemColors.Window;
            this.ClientSize = new System.Drawing.Size(903, 450);
            this.Controls.Add(this.statusText);
            this.Controls.Add(this.status);
            this.Controls.Add(this.progressBar1);
            this.Controls.Add(this.label4);
            this.Controls.Add(this.closeBtn);
            this.Controls.Add(this.label3);
            this.Controls.Add(this.pictureBox1);
            this.Controls.Add(this.submitTfa);
            this.Controls.Add(this.label2);
            this.Controls.Add(this.label1);
            this.Controls.Add(this.groupBox1);
            this.MaximizeBox = false;
            this.MinimizeBox = false;
            this.Name = "Form1";
            this.ShowIcon = false;
            this.StartPosition = System.Windows.Forms.FormStartPosition.CenterScreen;
            this.Text = "TRASA Windows Protection";
            this.TopMost = true;
            this.Load += new System.EventHandler(this.Form1_Load);
            ((System.ComponentModel.ISupportInitialize)(this.pictureBox1)).EndInit();
            this.groupBox1.ResumeLayout(false);
            this.groupBox1.PerformLayout();
            this.ResumeLayout(false);
            this.PerformLayout();

        }

        #endregion

        private System.Windows.Forms.Label label1;
        private System.Windows.Forms.Label label2;
        private System.Windows.Forms.TextBox totpCode;
        private System.Windows.Forms.Button submitTfa;
        private System.Windows.Forms.PictureBox pictureBox1;
        private System.Windows.Forms.GroupBox groupBox1;
        private System.Windows.Forms.Label label3;
        private System.Windows.Forms.ComboBox selectTfaMethod;
        private System.Windows.Forms.Button closeBtn;
        private System.Windows.Forms.TextBox trasaEmail;
        private System.Windows.Forms.Label label4;
        private System.Windows.Forms.ProgressBar progressBar1;
        private System.Windows.Forms.Label status;
        private System.Windows.Forms.Label statusText;
    }
}

