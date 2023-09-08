import http from 'k6/http';
import { describe, expect } from 'https://jslib.k6.io/k6chaijs/4.3.4.2/index.js';

export const options = {
  vus: 2,
  iterations: 2,
};

export default function () {
  describe('upsert', () => {
    const url = `${__ENV.MY_HOSTNAME}/api/channel/upsert-with-lock`;
    const payload = JSON.stringify([
      {
        code: makeCode(5),
      }
    ]);
  
    const params = {
      headers: {
        'Content-Type': 'application/json',
      },
    };
  
    const response = http.post(url, payload, params);
    console.log(response.body)

    expect(response.status, 'response status').to.equal(201);
    expect(response).to.have.validJsonBody();
  });
}

function makeCode(length) {
  let result = '';
  const characters = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ';
  const charactersLength = characters.length;
  let counter = 0;
  while (counter < length) {
    result += characters.charAt(Math.floor(Math.random() * charactersLength));
    counter += 1;
  }
  return result;
}