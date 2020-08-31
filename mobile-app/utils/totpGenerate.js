import JsSHA from 'jssha';

function hex2dec(s) {
  return parseInt(s, 16);
}

function dec2hex(s) {
  return (s < 15.5 ? '0' : '') + Math.round(s).toString(16);
}

function leftpad(str, len, pad) {
  let paddedStr;

  if (len + 1 >= str.length) {
    paddedStr = Array(len + 1 - str.length).join(pad) + str;
  }
  return paddedStr;
}

function base32tohex(base32) {
  const base32chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ234567';
  let bits = '';
  let hex = '';

  for (let i = 0; i < base32.length; i++) {
    const val = base32chars.indexOf(base32.charAt(i).toUpperCase());
    bits += leftpad(val.toString(2), 5, '0');
  }

  for (let i = 0; i + 4 <= bits.length; i += 4) {
    const chunk = bits.substr(i, 4);
    hex += parseInt(chunk, 2).toString(16);
  }
  return hex;
}

export default function getToken(key) {
  let decodedkey;
  try {
    // let epoch, time, shaObj, hmac, offset, otp;
    decodedkey = base32tohex(key);
    const epoch = Math.round(Date.now() / 1000.0);
    const time = leftpad(dec2hex(Math.floor(epoch / 30)), 16, '0');
    const shaObj = new JsSHA('SHA-1', 'HEX');
    shaObj.setHMACKey(decodedkey, 'HEX');
    shaObj.update(time);
    const hmac = shaObj.getHMAC('HEX');
    const offset = hex2dec(hmac.substring(hmac.length - 1));
    let otp = `${hex2dec(hmac.substr(offset * 2, 8)) & hex2dec('7fffffff')}`;
    otp = otp.substr(otp.length - 6, 6);
    return otp;
  } catch (e) {
    return 'Invalid secret';
  }
}
