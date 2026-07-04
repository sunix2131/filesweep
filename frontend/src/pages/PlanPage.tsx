import { useState } from 'react';
import { api } from '../lib/api';
import { bytes } from '../lib/format';
import { useAppStore } from '../stores/appStore';

export function PlanPage() {
  const { t, actionItems, clearPlan } = useAppStore();
  const [target, setTarget] = useState('');
  const total = actionItems.reduce((sum, i) => sum + i.sourceSizeBytes, 0);
  const execute = async (type: string) => {
    const items = actionItems.map((i) => ({ ...i, targetPath: target ? `${target}/${i.sourcePath.split(/[\\/]/).pop()}` : i.targetPath }));
    await api.execute(type, items); clearPlan();
  };
  return <section><header className="pageHeader"><h1>{t('plan')}</h1><button onClick={clearPlan}>{t('clear')}</button></header>
    <div className="toolbar"><input value={target} onChange={(e) => setTarget(e.target.value)} placeholder={t('folder')} /><button onClick={() => execute('move_to_folder')}>{t('move')}</button><button onClick={() => execute('move_to_trash')}>{t('trash')}</button></div>
    <div className="stats"><div className="stat"><span>{t('files')}</span><strong>{actionItems.length}</strong></div><div className="stat"><span>{t('size')}</span><strong>{bytes(total)}</strong></div></div>
    <table className="table"><tbody>{actionItems.map((i) => <tr key={i.sourcePath}><td>{i.sourcePath}</td><td>{bytes(i.sourceSizeBytes)}</td></tr>)}</tbody></table>
  </section>;
}
