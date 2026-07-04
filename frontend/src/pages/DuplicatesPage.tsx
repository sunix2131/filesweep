import { useEffect, useState } from 'react';
import { api } from '../lib/api';
import { bytes } from '../lib/format';
import { useAppStore } from '../stores/appStore';
import type { DuplicateGroup } from '../types/filesweep';
import { FileTable } from '../components/FileTable';
import { Empty } from '../components/Empty';

export function DuplicatesPage() {
  const { t, activeScan } = useAppStore();
  const [groups, setGroups] = useState<DuplicateGroup[]>([]);
  const [open, setOpen] = useState<DuplicateGroup>();
  useEffect(() => { if (activeScan) api.listDuplicates(activeScan.id).then((r) => setGroups(r.items)); }, [activeScan]);
  const load = async (id: string) => setOpen(await api.getDuplicate(id));
  return <section><header className="pageHeader"><h1>{t('duplicates')}</h1></header>
    {groups.length ? <table className="table"><thead><tr><th>{t('copies')}</th><th>{t('size')}</th><th>{t('reclaimable')}</th><th>{t('category')}</th></tr></thead>
      <tbody>{groups.map((g) => <tr key={g.id} onClick={() => load(g.id)}><td>{g.filesCount}</td><td>{bytes(g.fileSizeBytes)}</td><td>{bytes(g.estimatedReclaimableBytes)}</td><td>{g.category}</td></tr>)}</tbody></table> : <Empty />}
    {open?.members && <div className="drawer"><h2>{t('recommended')}</h2><FileTable files={open.members} /></div>}
  </section>;
}
