'use strict';

import JsSHA from 'jssha';

export default function getToken(key) {
  try {
    let epoch, time, shaObj, hmac, offset, otp;
    key = base32tohex(key);
    epoch = Math.round(Date.now() / 1000.0);
    time = leftpad(dec2hex(Math.floor(epoch / 30)), 16, '0');
    shaObj = new JsSHA('SHA-1', 'HEX');
    shaObj.setHMACKey(key, 'HEX');
    shaObj.update(time);
    hmac = shaObj.getHMAC('HEX');
    offset = hex2dec(hmac.substring(hmac.length - 1));
    otp = (hex2dec(hmac.substr(offset * 2, 8)) & hex2dec('7fffffff')) + '';
    otp = otp.substr(otp.length - 6, 6);
    return otp;
  } catch {
    return 'Invalid secret';
  }
}

function hex2dec(s) {
  return parseInt(s, 16);
}

function dec2hex(s) {
  return (s < 15.5 ? '0' : '') + Math.round(s).toString(16);
}

function base32tohex(base32) {
  let base32chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ234567',
    bits = '',
    hex = '';

  for (let i = 0; i < base32.length; i++) {
    let val = base32chars.indexOf(base32.charAt(i).toUpperCase());
    bits += leftpad(val.toString(2), 5, '0');
  }

  for (let i = 0; i + 4 <= bits.length; i += 4) {
    let chunk = bits.substr(i, 4);
    hex = hex + parseInt(chunk, 2).toString(16);
  }
  return hex;
}

function leftpad(str, len, pad) {
  if (len + 1 >= str.length) {
    str = Array(len + 1 - str.length).join(pad) + str;
  }
  return str;
}
