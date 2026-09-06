import { useEffect, useState } from 'react';
import { api } from '../lib/api';
import { useAppStore } from '../stores/appStore';
import type { FileAction } from '../types/filesweep';
import { Empty } from '../components/Empty';
import { Notice } from '../components/Notice';
import { errorMessage } from '../lib/errors';

export function HistoryPage() {
  const t = useAppStore((s) => s.t);
  const [items, setItems] = useState<FileAction[]>([]);
  const [error, setError] = useState('');
  const [busyId, setBusyId] = useState('');
  const refresh = () => api.history().then((result) => { setItems(result); setError(''); });
  useEffect(() => {
    void refresh().catch((reason) => setError(errorMessage(reason)));
  }, []);
  const undo = async (id: string) => {
    setBusyId(id);
    try {
      await api.undo(id);
      await refresh();
    } catch (reason) {
      setError(errorMessage(reason));
    } finally {
      setBusyId('');
    }
  };
  return <section><header className="pageHeader"><h1>{t('history')}</h1></header>
    {error && <Notice>{error}</Notice>}
    {items.length ? <table className="table"><tbody>{items.map((a) => <tr key={a.id}><td>{a.summary}</td><td>{a.status}</td><td>{a.undoAvailable && <button disabled={busyId === a.id} onClick={() => void undo(a.id)}>{t('undo')}</button>}</td></tr>)}</tbody></table> : <Empty />}
  </section>;
}
