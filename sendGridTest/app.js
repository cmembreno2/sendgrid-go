const sgMail= require('@sendgrid/mail');
const apiKey = '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';

sgMail.setApiKey(apiKey);

const message={
    to:'cmembreno@getmaya.com',
    from: {
        name:'SIEMPRE REI',
        email:'siemprerei@gmail.com'
    },
    subject: 'Sending with SendGrid is Fun',
    text: 'and easy to do anywhere, even with Node.js',
    html: '<strong>and easy to do anywhere, even with Node.js</strong>'
};

sgMail.send(message).then((response) => console.log('Email sent...')).catch((error)=>console.log(error.message));