import { Button, TextInput, Text, NumberInput } from '@mantine/core';
import { useMutation, useQuery } from '@tanstack/react-query';
import { useState } from 'react';
import './App.css';
import { TableSelection } from './components/TableSelection';
import { Applications, Config } from './types';

const apiUrl = 'http://localhost:6970';

const getApps = () => fetch(apiUrl + '/api/applications')
  .then(x => x.json())
  .then((apps: Applications) => apps.items.map(x => x.metadata.name));

const getConfig = () => fetch(apiUrl + '/api/config')
  .then(x => x.json())
  .then(x => x as Config);

const saveConfig = (data: Config) => fetch(apiUrl + '/api/config', {
  method: 'POST',
  body: JSON.stringify(data)
});

function App() {
  const [selection, setSelection] = useState<string[]>([]);
  const [envSelector, setEnvSelector] = useState('');
  const [refreshInterval, setRefreshInterval] = useState(60);

  const { data: apps } = useQuery(['apps'], getApps);
  useQuery(['config'], getConfig, {
    onSuccess: config => {
      setSelection(config?.allowed ?? []);
      setEnvSelector(config?.selectors?.environment ?? '');
      setRefreshInterval(config?.refreshInterval ?? 60);
    },
  });

  const { mutate, isLoading } = useMutation(saveConfig, {
    onError: err => alert(`Failed to save: ${err}`)
  });

  return (
    <div className="App">
      <TextInput
        placeholder="prod"
        label="Environment selector"
        value={envSelector}
        onChange={e => setEnvSelector(e.target.value)}
      />

      <NumberInput
        label="Refresh interval in seconds"
        value={refreshInterval}
        onChange={x => setRefreshInterval(x!)}
        required
      />
      <Text size="sm" pt="sm">
        Applications to sync
      </Text>
      <TableSelection data={apps || []} onSelectionChange={setSelection} selection={selection} />

      <Button 
        onClick={() => mutate({ 
          allowed: selection,
          refreshInterval: refreshInterval,
          selectors: {
            environment: envSelector,
          }
        })} 
        loading={isLoading}
        mt="sm"
        fullWidth
      >
        Save
      </Button>
    </div>
  );
}

export default App;
