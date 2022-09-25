import { Applications, Config } from './types';

const apiUrl = 'http://localhost:6970';

export const getApps = () => fetch(apiUrl + '/api/applications')
  .then(x => x.json())
  .then((apps: Applications) => apps.items.map(x => x.metadata.name));

export const getConfig = () => fetch(apiUrl + '/api/config')
  .then(x => x.json())
  .then(x => x as Config);

export const saveConfig = (data: Config) => fetch(apiUrl + '/api/config', {
  method: 'POST',
  body: JSON.stringify(data)
});
