const data = require('./data.json');
require('fs').writeFileSync('country', data.country);
