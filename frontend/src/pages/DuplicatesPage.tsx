import { useEffect, useState } from 'react';
import { api } from '../lib/api';
import { bytes } from '../lib/format';
import { useAppStore } from '../stores/appStore';
import type { DuplicateGroup } from '../types/filesweep';
import { FileTable } from '../components/FileTable';
import { Empty } from '../components/Empty';
import { Notice } from '../components/Notice';
import { errorMessage } from '../lib/errors';

export function DuplicatesPage() {
  const { t, activeScan } = useAppStore();
  const [groups, setGroups] = useState<DuplicateGroup[]>([]);
  const [open, setOpen] = useState<DuplicateGroup>();
  const [error, setError] = useState('');
  useEffect(() => {
    if (!activeScan) {
      setGroups([]);
      setOpen(undefined);
      return;
    }
    let active = true;
    void api.listDuplicates(activeScan.id)
      .then((result) => { if (active) { setGroups(result.items); setError(''); } })
      .catch((reason) => { if (active) setError(errorMessage(reason)); });
    return () => { active = false; };
  }, [activeScan]);
  const load = async (id: string) => {
    try {
      setOpen(await api.getDuplicate(id));
      setError('');
    } catch (reason) {
      setError(errorMessage(reason));
    }
  };
  return <section><header className="pageHeader"><h1>{t('duplicates')}</h1></header>
    {error && <Notice>{error}</Notice>}
    {groups.length ? <table className="table"><thead><tr><th>{t('copies')}</th><th>{t('size')}</th><th>{t('reclaimable')}</th><th>{t('category')}</th></tr></thead>
      <tbody>{groups.map((g) => <tr className="interactiveRow" key={g.id} tabIndex={0} onClick={() => void load(g.id)} onKeyDown={(event) => { if (event.key === 'Enter' || event.key === ' ') void load(g.id); }}><td>{g.filesCount}</td><td>{bytes(g.fileSizeBytes)}</td><td>{bytes(g.estimatedReclaimableBytes)}</td><td>{g.category}</td></tr>)}</tbody></table> : <Empty />}
    {open?.members && <div className="drawer"><h2>{t('recommended')}</h2><FileTable files={open.members} /></div>}
  </section>;
}
