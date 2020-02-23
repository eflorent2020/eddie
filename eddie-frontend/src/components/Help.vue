<template>
  <v-card>
    <v-card-title>
    <h1>Help</h1>
      </v-card-title>        
    <v-card-text>
      <h2>API</h2>
      <h3>API Document</h3>
      <p>
        <code>POST /document/:document_id</code>
        <h3>Params</h3>
        <ul>
          <li>required <var>email</var> your email</li>          
          <li>required <var>api_key</var> your api key</li>
          <li>optional <var>data</var> template data as json</li>          
          <li>optional <var>context</var> test or request_signature</li>       
          
          <li>optional <var>signature</var> signature reference</li>
          <li>optional <var>pincode</var> password to unlock your P12 file</li>
          <li>optional <var>sign_reason</var> text</li>
          <li>optional <var>sign_location</var> text</li>
          <li>optional <var>sign_contact</var> id</li>
          
        </ul>
      </p> 
      <h3>API Mail</h3>
      <p>
        <code>POST /mail/:mail_id</code>
        <h3>Params</h3>
        <ul>
          <li>required <var>email</var> your email</li>          
          <li>required <var>api_key</var> your api key</li>
          <li>optional <var>data</var> template data as json</li>          
          <li>optional <var>context</var> test or request_signature</li>       
          
          <li>optional <var>from</var> the From: email header field</li>
          <li>optional <var>to</var> the From: email header field</li>
          <li>optional <var>attachments</var> an array of pdf template id</li>

          <li>optional <var>signature</var> signature reference</li>
          <li>optional <var>pincode</var> password to unlock your P12 file</li>
          <li>optional <var>sign_reason</var> text</li>
          <li>optional <var>sign_location</var> text</li>
          <li>optional <var>sign_contact</var> id</li>
          
        </ul>
      </p>            
      <h2>Using the template engine</h2>
      <ul>
        <li><h3>Simple replace</h3>          
          <code>&lcub;&lcub;</code><code>.To</code><code> &rcub;&rcub;</code>
        </li> 
        <li><h3>Loops</h3>
         <code>&lcub;&lcub;</code> <code>range .Products</code>  <code>&rcub;&rcub;</code><br>
         <code>   &lcub;&lcub;</code> <code>.price</code>  <code>&rcub;&rcub;</code><br>
         <code>&lcub;&lcub;</code> <code>end</code>  <code>&rcub;&rcub;</code><br>
</li>
<li><h3>Conditionals</h3>

<code>&lcub;&lcub;</code><code>if .IsAttachment</code><code>&rcub;&rcub;</code><br>
<code>    -- attachment presentation code --</code><br>
<code>&lcub;&lcub;</code><code>end</code><code>&rcub;&rcub;</code>


        </li>
      </ul>
       
 <p></p>  
      <h2>QRCode</h2>      
      <p>In Templates simply set an image with the attribute src= starrting with Qrcode_XXXXX example:<br>
        <code>img src="&lcub;&lcub;Qrcode_isbn ... width=64 height=63</code><br>
      </p>
 <p></p>
      <h2>Getting a P12File</h2>
      <h3>Simple Electronic Signatures based on your domain name</h3>
          <p>You can use the certificate of your domain name (.com, .org, etc...) sometimes recognized as proof of identity.<br>
        To generate a .p12 from your cert use the following command:<br>
        <code>openssl pkcs12 -export -out certificate.p12 -inkey domain.key -in chained.pem</code>
      <h3>Using Letsencrypt for Simple Electronic signature</h3>
      <p>Letsencrypt offer free SSL/TLS certificate based the fact that you control your domain name
        We aim to provide Letsencrypt automatic certificate generation for you soon.
      </p>

      <h3>Advanced Electronic Signatures (AdES) & Qualified Electronic Signatures (QES)</h3>
      <p>This kind of signature requires you to sign the document offline using a certificate delivered 
        on a USB Key or Smartcard only, that you retain in your possession, and wich is delivered with identity validation.
        You'll then have to sign the document on your own device using your protected private key (smartcard or USB) 
        and then upload it back on our workflow.<br> 
      </p>
      <h3> </h3>


        </li>
      </ul>
    </v-card-text>
  </v-card>
</template>

<script>
export default {
  name: 'Help'
}
</script>
<style>
 ul { list-style-type: none; }
 li { 
   list-style-type: none; 
   margin-left: 24px;
   }

</style>
