import type { ScanFile } from '../types/filesweep';
import { bytes, date } from '../lib/format';
import { api } from '../lib/api';
import { useAppStore } from '../stores/appStore';
import { useState } from 'react';
import { Notice } from './Notice';
import { errorMessage } from '../lib/errors';

export function FileTable({ files }: { files: ScanFile[] }) {
  const { t, addAction } = useAppStore();
  const [error, setError] = useState('');
  const run = async (operation: () => Promise<unknown>) => {
    try {
      await operation();
      setError('');
    } catch (reason) {
      setError(errorMessage(reason));
    }
  };
  return (
    <>
    {error && <Notice>{error}</Notice>}
    <table className="table">
      <thead><tr><th>{t('name')}</th><th>{t('size')}</th><th>{t('category')}</th><th>{t('modified')}</th><th>{t('actions')}</th></tr></thead>
      <tbody>{files.map((f) => <tr key={f.id}>
        <td><span className="fileName">{f.name}</span><small>{f.absolutePath}</small></td>
        <td>{bytes(f.sizeBytes)}</td><td>{f.category}</td><td>{date(f.modifiedAt)}</td>
        <td className="rowActions">
          <button onClick={() => void run(() => api.open(f.absolutePath))}>{t('open')}</button>
          <button onClick={() => void run(() => api.reveal(f.absolutePath))}>{t('reveal')}</button>
          <button onClick={() => addAction({ sourcePath: f.absolutePath, sourceSizeBytes: f.sizeBytes, sourceSha256: f.sha256 })}>{t('addToPlan')}</button>
        </td>
      </tr>)}</tbody>
    </table>
    </>
  );
}
