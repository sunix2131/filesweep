import { useState } from 'react';
import { api } from '../lib/api';
import { bytes } from '../lib/format';
import { useAppStore } from '../stores/appStore';
import { Notice } from '../components/Notice';
import { errorMessage } from '../lib/errors';

export function PlanPage() {
  const { t, actionItems, clearPlan } = useAppStore();
  const [target, setTarget] = useState('');
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState('');
  const [completed, setCompleted] = useState('');
  const total = actionItems.reduce((sum, i) => sum + i.sourceSizeBytes, 0);
  const execute = async (type: string) => {
    const items = actionItems.map((i) => ({ ...i, targetPath: target ? `${target}/${i.sourcePath.split(/[\\/]/).pop()}` : i.targetPath }));
    setBusy(true);
    setError('');
    setCompleted('');
    try {
      const result = await api.execute(type, items);
      if (result.status === 'failed') throw new Error(result.errorMessage || result.summary);
      clearPlan();
      setCompleted(t('actionCompleted'));
    } catch (reason) {
      setError(errorMessage(reason));
    } finally {
      setBusy(false);
    }
  };
  return <section><header className="pageHeader"><h1>{t('plan')}</h1><button onClick={clearPlan}>{t('clear')}</button></header>
    <div className="toolbar"><input aria-label={t('folder')} value={target} onChange={(e) => setTarget(e.target.value)} placeholder={t('folder')} /><button disabled={busy || actionItems.length === 0 || !target.trim()} onClick={() => void execute('move_to_folder')}>{t('move')}</button><button disabled={busy || actionItems.length === 0} onClick={() => void execute('move_to_trash')}>{t('trash')}</button></div>
    {error && <Notice>{error}</Notice>}
    {completed && <Notice kind="success">{completed}</Notice>}
    <div className="stats"><div className="stat"><span>{t('files')}</span><strong>{actionItems.length}</strong></div><div className="stat"><span>{t('size')}</span><strong>{bytes(total)}</strong></div></div>
    <table className="table"><tbody>{actionItems.map((i) => <tr key={i.sourcePath}><td>{i.sourcePath}</td><td>{bytes(i.sourceSizeBytes)}</td></tr>)}</tbody></table>
  </section>;
}
