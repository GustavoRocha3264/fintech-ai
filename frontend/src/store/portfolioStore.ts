import { create } from 'zustand';
import type { AddPositionInput, Currency, Dashboard } from '../domain/models';
import { api } from '../services/api';

const STORAGE_KEY = 'cbpi.portfolioId';

function persistedID(): string | null {
  try {
    return localStorage.getItem(STORAGE_KEY);
  } catch {
    return null;
  }
}

function persist(id: string | null): void {
  try {
    if (id) localStorage.setItem(STORAGE_KEY, id);
    else localStorage.removeItem(STORAGE_KEY);
  } catch {
    /* ignore */
  }
}

interface PortfolioState {
  portfolioId: string | null;
  dashboard: Dashboard | null;
  loading: boolean;
  busy: boolean; // for in-flight mutations (add position, run analysis)
  error: string | null;

  createPortfolio: (baseCurrency: Currency) => Promise<void>;
  loadDashboard: (id?: string) => Promise<void>;
  addPosition: (input: AddPositionInput) => Promise<void>;
  runAnalysis: () => Promise<void>;
  reset: () => void;
}

export const usePortfolioStore = create<PortfolioState>((set, get) => ({
  portfolioId: persistedID(),
  dashboard: null,
  loading: false,
  busy: false,
  error: null,

  async createPortfolio(baseCurrency) {
    set({ busy: true, error: null });
    try {
      const portfolio = await api.createPortfolio({ baseCurrency });
      persist(portfolio.id);
      set({ portfolioId: portfolio.id, busy: false });
      await get().loadDashboard(portfolio.id);
    } catch (e) {
      set({ error: (e as Error).message, busy: false });
    }
  },

  async loadDashboard(id) {
    const targetId = id ?? get().portfolioId;
    if (!targetId) return;
    set({ loading: true, error: null });
    try {
      const dashboard = await api.getDashboard(targetId);
      set({ dashboard, loading: false });
    } catch (e) {
      set({ error: (e as Error).message, loading: false });
    }
  },

  async addPosition(input) {
    const id = get().portfolioId;
    if (!id) {
      set({ error: 'No portfolio selected' });
      return;
    }
    set({ busy: true, error: null });
    try {
      await api.addPosition(id, input);
      await get().loadDashboard(id);
      set({ busy: false });
    } catch (e) {
      set({ error: (e as Error).message, busy: false });
    }
  },

  async runAnalysis() {
    const id = get().portfolioId;
    if (!id) {
      set({ error: 'No portfolio selected' });
      return;
    }
    set({ busy: true, error: null });
    try {
      await api.runAnalysis(id);
      await get().loadDashboard(id);
      set({ busy: false });
    } catch (e) {
      set({ error: (e as Error).message, busy: false });
    }
  },

  reset() {
    persist(null);
    set({ portfolioId: null, dashboard: null, error: null });
  },
}));
