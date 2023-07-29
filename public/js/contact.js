function submitData() {
  let name = document.getElementById("input-name").value;
  let email = document.getElementById("input-email").value;
  let phone = document.getElementById("input-phone").value;
  let subject = document.getElementById("input-subject").value;
  let message = document.getElementById("input-message").value;

  if (name == "") {
    return alert("Please input your name !");
  } else if (email == "") {
    return alert("Please input your email !");
  } else if (phone == "") {
    return alert("Please input your Phone Number !");
  } else if (subject == "") {
    return alert("Please input your Subject !");
  } else if (message == "") {
    return alert("Please input your Message !");
  }

  let emailReceiver = "angga.ardiansyah056@gmail.com";

  let a = document.createElement("a");
  a.href = `mailto:${emailReceiver}?subject=${subject}&body=Halo, nama saya ${name}, ${message}. Silahkan kontak saya di nomor ${phone}, terimakasih.`;
  a.click();

  // console.log(name);
  // console.log(email);
  // console.log(phone);
  // console.log(subject);
  // console.log(message);

  let emailer = {
    name,
    email,
    phone,
    subject,
    message,
  };
  console.log(emailer);
}
