import { create } from 'zustand';
import ru from '../i18n/ru.json';
import en from '../i18n/en.json';
import type { ActionItem, ScanProgress, ScanSession, Settings } from '../types/filesweep';

const dictionaries = { ru, en } as const;
type Lang = keyof typeof dictionaries;

type State = {
  lang: Lang;
  theme: 'system' | 'light' | 'dark';
  selectedPaths: string[];
  activeScan?: ScanSession;
  scanProgress?: ScanProgress;
  actionItems: ActionItem[];
  t: (key: keyof typeof ru) => string;
  setSettings: (settings: Settings) => void;
  addPath: (path: string) => void;
  setActiveScan: (scan: ScanSession) => void;
  setScanProgress: (progress?: ScanProgress) => void;
  addAction: (item: ActionItem) => void;
  clearPlan: () => void;
};

export const useAppStore = create<State>((set, get) => ({
  lang: 'ru',
  theme: 'system',
  selectedPaths: [],
  actionItems: [],
  t: (key) => dictionaries[get().lang][key] ?? key,
  setSettings: (settings) => set({ lang: settings.language, theme: settings.theme }),
  addPath: (path) => set((s) => ({ selectedPaths: [...new Set([...s.selectedPaths, path])] })),
  setActiveScan: (scan) => set({ activeScan: scan }),
  setScanProgress: (progress) => set({ scanProgress: progress }),
  addAction: (item) => set((s) => ({ actionItems: [...s.actionItems, item] })),
  clearPlan: () => set({ actionItems: [] })
}));
