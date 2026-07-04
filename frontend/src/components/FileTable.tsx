import type { ScanFile } from '../types/filesweep';
import { bytes, date } from '../lib/format';
import { api } from '../lib/api';
import { useAppStore } from '../stores/appStore';

export function FileTable({ files }: { files: ScanFile[] }) {
  const { t, addAction } = useAppStore();
  return (
    <table className="table">
      <thead><tr><th>{t('name')}</th><th>{t('size')}</th><th>{t('category')}</th><th>{t('modified')}</th><th>{t('actions')}</th></tr></thead>
      <tbody>{files.map((f) => <tr key={f.id}>
        <td><span className="fileName">{f.name}</span><small>{f.absolutePath}</small></td>
        <td>{bytes(f.sizeBytes)}</td><td>{f.category}</td><td>{date(f.modifiedAt)}</td>
        <td className="rowActions">
          <button onClick={() => api.open(f.absolutePath)}>{t('open')}</button>
          <button onClick={() => api.reveal(f.absolutePath)}>{t('reveal')}</button>
          <button onClick={() => addAction({ sourcePath: f.absolutePath, sourceSizeBytes: f.sizeBytes, sourceSha256: f.sha256 })}>{t('addToPlan')}</button>
        </td>
      </tr>)}</tbody>
    </table>
  );
}
