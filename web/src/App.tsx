import { useQuery } from '@tanstack/react-query';
import './App.css';
import { TableSelection } from './components/TableSelection';

type Applications = {
  items: Array<{
    metadata: {
      name: string;
    }
  }>;
}

const apiUrl = 'http://localhost:6970';
const getApps = () => fetch(apiUrl + '/api/applications')
  .then(x => x.json())
  .then((apps: Applications) => apps.items.map(x => x.metadata.name));

function App() {
  const { data } = useQuery(['apps'], getApps);

  return (
    <div className="App">
      <TableSelection data={data || []} />
    </div>
  );
}

export default App;
