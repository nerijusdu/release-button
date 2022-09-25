import { Button } from '@mantine/core';
import { useMutation, useQuery } from '@tanstack/react-query';
import { useState } from 'react';
import './App.css';
import { TableSelection } from './components/TableSelection';

type Applications = {
  items: Array<{
    metadata: { name: string }
  }>;
}

type SaveConfigRequest = {
  allowedApps: string[];
}

type Config = {
  allowed: string[];
}

const apiUrl = 'http://localhost:6970';

const getApps = () => fetch(apiUrl + '/api/applications')
  .then(x => x.json())
  .then((apps: Applications) => apps.items.map(x => x.metadata.name));

const getConfig = () => fetch(apiUrl + '/api/config')
  .then(x => x.json())
  .then(x => x as Config);

const saveConfig = (data: SaveConfigRequest) => fetch(apiUrl + '/api/config', {
  method: 'POST',
  body: JSON.stringify(data)
});

function App() {
  const [selection, setSelection] = useState<string[]>([]);
  const { data: apps } = useQuery(['apps'], getApps);
  useQuery(['config'], getConfig, {
    onSuccess: config => {
      setSelection(config?.allowed ?? []);
    }
  });

  const { mutate, isLoading } = useMutation(saveConfig, {
    onError: err => alert(`Failed to save: ${err}`)
  });

  return (
    <div className="App">
      <TableSelection data={apps || []} onSelectionChange={setSelection} selection={selection} />
      <Button 
        onClick={() => mutate({ allowedApps: selection })} 
        loading={isLoading}
        fullWidth
      >
        Save
      </Button>
    </div>
  );
}

export default App;
