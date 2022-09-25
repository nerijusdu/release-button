import { Button, TextInput, Text, NumberInput } from '@mantine/core';
import { useMutation, useQuery } from '@tanstack/react-query';
import { useState } from 'react';
import { TableSelection } from './TableSelection';
import { getApps, getConfig, saveConfig } from './api';
import './App.css';

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
