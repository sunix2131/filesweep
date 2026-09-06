import { useEffect, useState } from 'react';
import { api } from '../lib/api';
import { bytes, date } from '../lib/format';
import { useAppStore } from '../stores/appStore';
import { Stat } from '../components/Stat';
import { Notice } from '../components/Notice';
import { errorMessage } from '../lib/errors';
import type { ScanSession } from '../types/filesweep';

export function HomePage() {
  const { t, selectedPaths, addPath, activeScan, setActiveScan, scanProgress, setScanProgress } = useAppStore();
  const [scans, setScans] = useState<ScanSession[]>([]);
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState('');
  useEffect(() => {
    let active = true;
    void api.listScans().then((items) => { if (active) setScans(items); }).catch((reason) => { if (active) setError(errorMessage(reason)); });
    return () => { active = false; };
  }, [activeScan]);
  const select = async () => {
    setError('');
    try {
      const res = await api.selectFolders();
      res.paths.forEach(addPath);
    } catch (reason) {
      setError(errorMessage(reason));
    }
  };
  const scan = async () => {
    setBusy(true);
    setError('');
    try {
      const started = await api.startScan(selectedPaths);
      setActiveScan(started);
      setScanProgress({ scanId: started.id, phase: 'queued', processedFiles: 0, totalFiles: 0, currentPath: '', percent: 0, status: 'running' });
    } catch (reason) {
      setError(errorMessage(reason));
    } finally {
      setBusy(false);
    }
  };
  const cancel = async () => {
    if (scanProgress?.scanId) {
      try {
        await api.cancelScan(scanProgress.scanId);
      } catch (reason) {
        setError(errorMessage(reason));
      }
    }
  };
  const latest = activeScan ?? scans[0];
  const showProgress = busy || scanProgress?.status === 'running';
  const percent = Math.max(0, Math.min(100, scanProgress?.percent ?? 0));
  return <section>
    <header className="pageHeader"><h1>{t('home')}</h1><button onClick={select}>{t('scanFolder')}</button></header>
    {error && <Notice>{error}</Notice>}
    <div className="toolbar"><button onClick={select}>{t('addFolder')}</button><button disabled={busy || selectedPaths.length === 0} onClick={scan}>{t('startScan')}</button></div>
    <div className="pathList">{selectedPaths.map((p) => <div key={p}>{p}</div>)}</div>
    {showProgress && <div className="scanProgress">
      <div className="scanProgressTop">
        <strong>{t('scanInProgress')}</strong>
        <span>{t(phaseLabelKey(scanProgress?.phase ?? 'queued'))}</span>
      </div>
      <div className={`progress ${scanProgress?.totalFiles ? '' : 'progressIndeterminate'}`}>
        <div style={scanProgress?.totalFiles ? { width: `${percent}%` } : undefined} />
      </div>
      <div className="scanProgressMeta">
        <span>{scanProgress?.totalFiles ? `${scanProgress.processedFiles} / ${scanProgress.totalFiles}` : `${t('foundFiles')}: ${scanProgress?.processedFiles ?? latest?.filesCount ?? 0}`}</span>
        <span>{scanProgress?.totalFiles ? `${Math.round(percent)}%` : t('preparing')}</span>
      </div>
      <div className="currentFile">
        <span>{t('currentFile')}</span>
        <code>{scanProgress?.currentPath || selectedPaths[0] || t('preparing')}</code>
      </div>
      <button onClick={cancel}>{t('cancel')}</button>
    </div>}
    {scanProgress?.status === 'completed' && <div className="scanNotice">{t('scanCompleted')}</div>}
    {scanProgress?.status === 'cancelled' && <div className="scanNotice">{t('scanCancelled')}</div>}
    {scanProgress?.status === 'failed' && <div className="scanNotice scanNoticeError">{scanProgress.message || t('scanFailed')}</div>}
    {latest && <div className="stats">
      <Stat label={t('files')} value={latest.filesCount} /><Stat label={t('size')} value={bytes(latest.totalSizeBytes)} />
      <Stat label={t('duplicateGroups')} value={latest.duplicateGroupsCount} /><Stat label={t('reclaimable')} value={bytes(latest.reclaimableBytes)} />
      <Stat label={t('errors')} value={latest.errorsCount} /><Stat label={t('latestScan')} value={date(latest.completedAt ?? latest.startedAt)} />
    </div>}
  </section>;
}

function phaseLabelKey(phase: string) {
  const labels: Record<string, 'phaseQueued' | 'phaseDiscovering' | 'phaseAnalyzing' | 'phaseHashing' | 'phaseSaving' | 'phaseCompleted' | 'phaseCancelled' | 'phaseFailed'> = {
    queued: 'phaseQueued',
    discovering: 'phaseDiscovering',
    analyzing: 'phaseAnalyzing',
    hashing: 'phaseHashing',
    saving: 'phaseSaving',
    completed: 'phaseCompleted',
    cancelled: 'phaseCancelled',
    failed: 'phaseFailed'
  };
  return labels[phase] ?? 'phaseQueued';
}
