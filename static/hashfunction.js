async function sha256(message) {
    // Convert the string to an ArrayBuffer
    const encoder = new TextEncoder();
    const data = encoder.encode(message);
  
    // Calculate the SHA-256 hash
    const hashBuffer = await crypto.subtle.digest('SHA-256', data);
  
    // Convert the hash to a hexadecimal string
    const hashArray = Array.from(new Uint8Array(hashBuffer));
    const hashHex = hashArray.map(byte => byte.toString(16).padStart(2, '0')).join('');
  
    return hashHex;
  }
  
  // Example usage
  const element1 = 'Hello';
  const element2 = ' ';
  const element3 = 'World';
  
  const concatenatedString = element1 + element2 + element3;
  
  sha256(concatenatedString)
    .then(hash => {
      console.log('SHA-256 hash:', hash);
    })
    .catch(error => {
      console.error('Error calculating hash:', error);
    });