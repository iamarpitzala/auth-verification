// Google Apps Script code to deploy at the webhook URL
// This should be deployed as a web app with execute permissions set to "Anyone"

function doPost(e) {
  try {
    // Parse the JSON payload
    const data = JSON.parse(e.postData.contents);
    const toEmail = data.toEmail;
    const verificationCode = data.verificationCode;
    
    if (!toEmail || !verificationCode) {
      return ContentService
        .createTextOutput(JSON.stringify({error: "Missing required fields"}))
        .setMimeType(ContentService.MimeType.JSON);
    }
    
    // Send verification email
    const result = sendVerificationEmail(toEmail, verificationCode);
    
    return ContentService
      .createTextOutput(JSON.stringify({success: true, code: result}))
      .setMimeType(ContentService.MimeType.JSON);
      
  } catch (error) {
    console.error('Error in doPost:', error);
    return ContentService
      .createTextOutput(JSON.stringify({error: error.toString()}))
      .setMimeType(ContentService.MimeType.JSON);
  }
}

function sendVerificationEmail(toEmail, verificationCode) {
  // Email subject
  const subject = "Your Verification Code";
  
  // Email body with the provided verification code
  const body = `Hello,

Your verification code is: ${verificationCode}

This code is valid for 10 minutes.

If you did not request this, please ignore this email.

Regards,
Support Team`;
  
  // Send email
  MailApp.sendEmail({
    to: toEmail,
    subject: subject,
    body: body
  });
  
  // Return the code for confirmation
  return verificationCode;
}

// Test function (optional)
function testSendEmail() {
  const result = sendVerificationEmail("test@example.com", "123456");
  console.log("Test email sent with code:", result);
}