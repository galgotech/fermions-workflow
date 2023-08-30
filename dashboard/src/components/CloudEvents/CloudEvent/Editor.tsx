import MonacoEditor from '@monaco-editor/react';

interface Props {
  value: string;
  onChange: (value: string) => void;
};

export const Editor = ({ value, onChange }: Props) => {
  return (
    <MonacoEditor
      width="100%"
      height="100%"
      language="json"
      value={value}
      options={{
        selectOnLineNumbers: true,
      }}
      onChange={(value) => {
        onChange(value || '');
      }}
      onMount={(editor) => {
        editor.focus();
      }}
    />
  );
};
