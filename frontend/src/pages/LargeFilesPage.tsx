import { useEffect, useState } from 'react';
import { api } from '../lib/api';
import { useAppStore } from '../stores/appStore';
import type { ScanFile } from '../types/filesweep';
import { FileTable } from '../components/FileTable';
import { Empty } from '../components/Empty';

export function LargeFilesPage() {
  const { t, activeScan } = useAppStore();
  const [files, setFiles] = useState<ScanFile[]>([]);
  const [min, setMin] = useState(100 * 1024 * 1024);
  const [search, setSearch] = useState('');
  useEffect(() => { if (activeScan) api.listFiles(activeScan.id, min, '', search).then((r) => setFiles(r.items)); }, [activeScan, min, search]);
  return <section><header className="pageHeader"><h1>{t('large')}</h1></header>
    <div className="toolbar"><input placeholder={t('search')} value={search} onChange={(e) => setSearch(e.target.value)} /><select value={min} onChange={(e) => setMin(Number(e.target.value))}><option value={104857600}>100 MB</option><option value={524288000}>500 MB</option><option value={1073741824}>1 GB</option></select></div>
    {files.length ? <FileTable files={files} /> : <Empty />}
  </section>;
}
