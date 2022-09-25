import { createStyles, Table, Checkbox, ScrollArea } from '@mantine/core';

const useStyles = createStyles((theme) => ({
  rowSelected: {
    backgroundColor:
      theme.colorScheme === 'dark'
        ? theme.fn.rgba(theme.colors[theme.primaryColor][7], 0.2)
        : theme.colors[theme.primaryColor][0],
  },
}));

interface TableSelectionProps {
  data: string[];
  onSelectionChange: (selection: string[]) => void;
  selection: string[];
}

export function TableSelection({ data, onSelectionChange, selection }: TableSelectionProps) {
  const { classes, cx } = useStyles();
  const toggleRow = (id: string) =>
    onSelectionChange(
      selection.includes(id) ? selection.filter((item) => item !== id) : [...selection, id]
    );
  const toggleAll = () => onSelectionChange(selection.length === data.length ? [] : data);

  const rows = data.map((item) => {
    const selected = selection.includes(item);
    return (
      <tr key={item} className={cx({ [classes.rowSelected]: selected })}>
        <td>
          <Checkbox
            checked={selection.includes(item)}
            onChange={() => toggleRow(item)}
            transitionDuration={0}
          />
        </td>
        <td>{item}</td>
      </tr>
    );
  });

  return (
    <ScrollArea>
      <Table verticalSpacing="sm">
        <thead>
          <tr>
            <th style={{ width: 40 }}>
              <Checkbox
                onChange={toggleAll}
                checked={selection.length === data.length && data.length}
                indeterminate={selection.length > 0 && selection.length !== data.length}
                transitionDuration={0}
              />
            </th>
            <th>Name</th>
          </tr>
        </thead>
        <tbody>{rows}</tbody>
      </Table>
    </ScrollArea>
  );
}
