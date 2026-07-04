import { useEffect, useState } from 'react';
import { api } from '../lib/api';
import { useAppStore } from '../stores/appStore';
import type { FileAction } from '../types/filesweep';
import { Empty } from '../components/Empty';

export function HistoryPage() {
  const t = useAppStore((s) => s.t);
  const [items, setItems] = useState<FileAction[]>([]);
  const refresh = () => api.history().then(setItems);
  useEffect(() => {
    void refresh();
  }, []);
  return <section><header className="pageHeader"><h1>{t('history')}</h1></header>
    {items.length ? <table className="table"><tbody>{items.map((a) => <tr key={a.id}><td>{a.summary}</td><td>{a.status}</td><td>{a.undoAvailable && <button onClick={() => api.undo(a.id).then(refresh)}>{t('undo')}</button>}</td></tr>)}</tbody></table> : <Empty />}
  </section>;
}
