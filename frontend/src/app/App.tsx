import { useEffect } from 'react';
import { NavLink, Route, Routes } from 'react-router-dom';
import { Archive, BarChart3, FolderSearch, History, Home, ListChecks, Moon, Settings, Trash2 } from 'lucide-react';
import { api } from '../lib/api';
import { onWailsEvent } from '../lib/wailsEvents';
import { useAppStore } from '../stores/appStore';
import type { ScanProgress } from '../types/filesweep';
import { HomePage } from '../pages/HomePage';
import { DuplicatesPage } from '../pages/DuplicatesPage';
import { LargeFilesPage } from '../pages/LargeFilesPage';
import { CategoriesPage } from '../pages/CategoriesPage';
import { PlanPage } from '../pages/PlanPage';
import { HistoryPage } from '../pages/HistoryPage';
import { SettingsPage } from '../pages/SettingsPage';
import { ErrorBoundary } from './ErrorBoundary';

const nav = [
  ['/', 'home', Home], ['/duplicates', 'duplicates', FolderSearch], ['/large', 'large', Archive],
  ['/categories', 'categories', BarChart3], ['/plan', 'plan', ListChecks], ['/history', 'history', History], ['/settings', 'settings', Settings]
] as const;

export function App() {
  const { t, setSettings, theme, lang, setScanProgress, setActiveScan } = useAppStore();
  useEffect(() => { api.getSettings().then(setSettings).catch(() => undefined); }, [setSettings]);
  useEffect(() => {
    const offProgress = onWailsEvent<ScanProgress>('scan:progress', (progress) => {
      setScanProgress({ ...progress, status: 'running' });
    });
    const offCompleted = onWailsEvent<string>('scan:completed', (scanId) => {
      const current = useAppStore.getState().scanProgress;
      setScanProgress(current?.scanId === scanId
        ? { ...current, phase: 'completed', percent: 100, status: 'completed' }
        : { scanId, phase: 'completed', processedFiles: 0, totalFiles: 0, currentPath: '', percent: 100, status: 'completed' });
      api.getScan(scanId).then(setActiveScan).catch(() => undefined);
    });
    const offCancelled = onWailsEvent<string>('scan:cancelled', (scanId) => {
      setScanProgress({ scanId, phase: 'cancelled', processedFiles: 0, totalFiles: 0, currentPath: '', percent: 0, status: 'cancelled' });
    });
    const offFailed = onWailsEvent<{ scanId: string; message: string }>('scan:failed', (event) => {
      setScanProgress({ scanId: event.scanId, phase: 'failed', processedFiles: 0, totalFiles: 0, currentPath: '', percent: 0, status: 'failed', message: event.message });
    });
    return () => {
      offProgress();
      offCompleted();
      offCancelled();
      offFailed();
    };
  }, [setActiveScan, setScanProgress]);
  useEffect(() => {
    const media = window.matchMedia?.('(prefers-color-scheme: dark)');
    const applyTheme = () => document.documentElement.classList.toggle('dark', theme === 'dark' || (theme === 'system' && media?.matches));
    applyTheme();
    document.documentElement.lang = lang;
    if (theme === 'system') media?.addEventListener('change', applyTheme);
    return () => media?.removeEventListener('change', applyTheme);
  }, [theme, lang]);
  return (
    <ErrorBoundary>
    <div className="app">
      <aside className="sidebar">
        <div className="brand"><Trash2 size={22} /> <span>{t('app')}</span></div>
        <nav>{nav.map(([to, key, Icon]) => <NavLink key={to} to={to} end={to === '/'}><Icon size={18} /> {t(key)}</NavLink>)}</nav>
        <div className="sidebarFoot"><Moon size={16} /> <span>{t('version')}</span></div>
      </aside>
      <main className="main">
        <Routes>
          <Route path="/" element={<HomePage />} />
          <Route path="/duplicates" element={<DuplicatesPage />} />
          <Route path="/large" element={<LargeFilesPage />} />
          <Route path="/categories" element={<CategoriesPage />} />
          <Route path="/plan" element={<PlanPage />} />
          <Route path="/history" element={<HistoryPage />} />
          <Route path="/settings" element={<SettingsPage />} />
        </Routes>
      </main>
    </div>
    </ErrorBoundary>
  );
}
