import { create } from 'zustand';
import type { Portfolio, Report } from '../types/domain';
import { api } from '../services/api';

interface PortfolioState {
  current: Portfolio | null;
  report: Report | null;
  loading: boolean;
  error: string | null;
  load: (id: string) => Promise<void>;
  runAnalysis: (id: string) => Promise<void>;
}

export const usePortfolioStore = create<PortfolioState>((set) => ({
  current: null,
  report: null,
  loading: false,
  error: null,
  async load(id) {
    set({ loading: true, error: null });
    try {
      set({ current: await api.getPortfolio(id), loading: false });
    } catch (e) {
      set({ error: (e as Error).message, loading: false });
    }
  },
  async runAnalysis(id) {
    set({ loading: true, error: null });
    try {
      set({ report: await api.runAnalysis(id), loading: false });
    } catch (e) {
      set({ error: (e as Error).message, loading: false });
    }
  },
}));
