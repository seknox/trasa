const { Client } = require('pg')
// pools will use environment variables
// for connection information
const client = new Client({
    user: 'trasauser',
    host: 'localhost',
    database: 'trasadb',
    password: 'secretpassword',
    port: 3211,
  })
  
// clients will also use environment variables
// for connection information
const client = new Client()
await client.connect()
const res = await client.query('SELECT NOW()')
await client.end()