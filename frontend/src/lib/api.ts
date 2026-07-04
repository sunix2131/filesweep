import type { ActionItem, DuplicateGroup, FileAction, ScanFile, ScanSession, Settings } from '../types/filesweep';

type Backend = Record<string, (...args: unknown[]) => Promise<unknown>>;
const wails = (): Backend => (window as unknown as { go?: { wails?: { AppAPI?: Backend } } }).go?.wails?.AppAPI ?? mockApi;
const array = <T>(value: T[] | null | undefined): T[] => value ?? [];
const page = <T>(value: { items?: T[] | null; total?: number } | null | undefined): { items: T[]; total: number } => ({
  items: array(value?.items),
  total: value?.total ?? 0
});

export const api = {
  selectFolders: () => wails().SelectFolders() as Promise<{ paths: string[]; hasDangerousPath: boolean }>,
  startScan: (paths: string[], confirmDangerousFolders = false) => wails().StartScan({ paths, confirmDangerousFolders }) as Promise<ScanSession>,
  cancelScan: (scanId: string) => wails().CancelScan(scanId),
  getScan: (scanId: string) => wails().GetScanSession(scanId) as Promise<ScanSession>,
  listScans: async () => array(await wails().ListScanSessions() as ScanSession[] | null),
  listDuplicates: async (scanId: string) => page(await wails().ListDuplicateGroups({ scanId, limit: 100, offset: 0 }) as { items?: DuplicateGroup[] | null; total?: number } | null),
  getDuplicate: (id: string) => wails().GetDuplicateGroup(id) as Promise<DuplicateGroup>,
  listFiles: async (scanId: string, minSizeBytes = 0, category = '', search = '') => page(await wails().ListLargeFiles({ scanId, minSizeBytes, category, search, limit: 500, offset: 0 }) as { items?: ScanFile[] | null; total?: number } | null),
  categories: async (scanId: string) => array(await wails().ListCategories(scanId) as Array<{ category: string; filesCount: number; totalSizeBytes: number }> | null),
  execute: (actionType: string, items: ActionItem[]) => wails().ExecuteActionPlan({ actionType, items }) as Promise<FileAction>,
  history: async () => array(await wails().ListActionHistory({ limit: 100, offset: 0 }) as FileAction[] | null),
  undo: (actionId: string) => wails().UndoAction(actionId) as Promise<FileAction>,
  reveal: (path: string) => wails().RevealInFileManager(path),
  open: (path: string) => wails().OpenFile(path),
  exportCSV: (scanId: string) => wails().ExportScanReport(scanId, 'csv') as Promise<{ path: string }>,
  getSettings: () => wails().GetSettings() as Promise<Settings>,
  updateSettings: (settings: Settings) => wails().UpdateSettings(settings) as Promise<Settings>,
  clearThumbnails: () => wails().ClearThumbnailCache()
};

const mockApi = new Proxy({}, { get: () => async () => {
  throw new Error('FileSweep backend is available inside Wails. Run wails dev for native features.');
}}) as Backend;
