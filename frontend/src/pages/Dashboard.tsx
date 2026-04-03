import React, { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import api from '../services/api';
import { 
  Plus, 
  FileText, 
  Clock, 
  CheckCircle2, 
  XCircle, 
  Loader2, 
  Download,
  AlertCircle,
  RefreshCcw
} from 'lucide-react';
import { cn } from '../lib/utils';
import { useTranslation } from 'react-i18next';

interface Job {
  job_id: string;
  process_type: string;
  status: string;
  created_at: string;
  output_files?: string[];
}

const Dashboard: React.FC = () => {
  const { t, i18n } = useTranslation();
  const [jobs, setJobs] = useState<Job[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  const fetchJobs = async () => {
    try {
      const response = await api.get('/pdf/list');
      // Sort by created_at descending
      const sortedJobs = (response.data || []).sort((a: Job, b: Job) => 
        new Date(b.created_at).getTime() - new Date(a.created_at).getTime()
      );
      setJobs(sortedJobs);
    } catch (err) {
      setError(t('dashboard_page.error_fetch'));
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchJobs();
    // Poll for status updates every 10 seconds
    const interval = setInterval(() => {
      fetchJobs();
    }, 10000);
    
    return () => clearInterval(interval);
  }, []);

  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'completed': return <CheckCircle2 size={18} className="text-green-500" />;
      case 'failed': return <XCircle size={18} className="text-red-500" />;
      case 'processing': return <Loader2 size={18} className="text-blue-500 animate-spin" />;
      case 'awaiting_files': return <Clock size={18} className="text-yellow-500" />;
      default: return <Clock size={18} className="text-dark-400" />;
    }
  };

  const getStatusLabel = (status: string) => {
    switch (status) {
      case 'completed': return t('dashboard_page.status.completed');
      case 'failed': return t('dashboard_page.status.failed');
      case 'processing': return t('dashboard_page.status.processing');
      case 'awaiting_files': return t('dashboard_page.status.awaiting_files');
      case 'pending': return t('dashboard_page.status.pending');
      default: return status;
    }
  };

  const getTypeLabel = (type: string) => {
    switch (type) {
      case 'merge': return t('dashboard_page.types.merge');
      case 'split': return t('dashboard_page.types.split');
      case 'protect': return t('dashboard_page.types.protect');
      case 'unprotect': return t('dashboard_page.types.unprotect');
      case 'remove_password': return t('dashboard_page.types.remove_password');
      default: return type;
    }
  };

  const handleDownload = async (jobId: string, filename: string) => {
    try {
      const response = await api.get(`/pdf/download/${jobId}`, {
        params: { filename }
      });
      
      const { url } = response.data;
      if (url) {
        window.open(url, '_blank');
      }
    } catch (err) {
      console.error('Erro ao obter URL de download:', err);
      alert(t('dashboard_page.error_download'));
    }
  };

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="flex items-center justify-between mb-8">
        <div>
          <h1 className="text-3xl font-bold">{t('dashboard_page.title')}</h1>
          <p className="text-dark-400 mt-1">{t('dashboard_page.subtitle')}</p>
        </div>
        <div className="flex gap-4">
          <button 
            onClick={() => { setLoading(true); fetchJobs(); }}
            className="p-2.5 bg-dark-900 border border-dark-800 rounded-lg text-dark-300 hover:text-dark-50 transition-all"
            title={t('dashboard_page.refresh')}
          >
            <RefreshCcw size={20} className={loading ? "animate-spin" : ""} />
          </button>
          <Link
            to="/new-job"
            className="bg-blue-600 hover:bg-blue-700 text-white px-5 py-2.5 rounded-lg font-semibold flex items-center gap-2 transition-all shadow-lg shadow-blue-500/20"
          >
            <Plus size={20} />
            {t('dashboard_page.new_process')}
          </Link>
        </div>
      </div>

      {error && (
        <div className="bg-red-500/10 border border-red-500/20 text-red-500 p-4 rounded-xl flex items-center gap-3 mb-6">
          <AlertCircle size={20} />
          <span>{error}</span>
        </div>
      )}

      {loading && jobs.length === 0 ? (
        <div className="flex flex-col items-center justify-center py-20 text-dark-400 gap-4">
          <Loader2 size={40} className="animate-spin text-blue-500" />
          <p>{t('common.loading')}</p>
        </div>
      ) : jobs.length === 0 ? (
        <div className="bg-dark-900 border border-dark-800 rounded-2xl p-12 text-center">
          <div className="mx-auto h-16 w-16 bg-dark-800 rounded-2xl flex items-center justify-center text-dark-600 mb-4">
            <FileText size={32} />
          </div>
          <h3 className="text-xl font-semibold mb-2">{t('dashboard_page.no_jobs')}</h3>
          <p className="text-dark-400 mb-8 max-w-sm mx-auto">
            {t('dashboard_page.start_first')}
          </p>
          <Link
            to="/new-job"
            className="inline-flex bg-blue-600 hover:bg-blue-700 text-white px-6 py-3 rounded-xl font-bold transition-all"
          >
            {t('dashboard_page.new_process')}
          </Link>
        </div>
      ) : (
        <div className="grid gap-4">
          {jobs.map((job) => (
            <div 
              key={job.job_id}
              className="bg-dark-900 border border-dark-800 rounded-xl p-4 sm:p-6 flex flex-col sm:flex-row sm:items-center justify-between gap-4 hover:border-dark-700 transition-all group"
            >
              <div className="flex items-center gap-4">
                <div className={cn(
                  "h-12 w-12 rounded-lg flex items-center justify-center",
                  job.status === 'completed' ? "bg-green-500/10 text-green-500" : 
                  job.status === 'failed' ? "bg-red-500/10 text-red-500" : "bg-blue-500/10 text-blue-500"
                )}>
                  <FileText size={24} />
                </div>
                <div>
                  <div className="flex items-center gap-2">
                    <h3 className="font-bold text-lg leading-none capitalize">
                      {getTypeLabel(job.process_type)}
                    </h3>
                    <span className="text-xs text-dark-500 font-mono">#{job.job_id.slice(0, 8)}</span>
                  </div>
                  <p className="text-sm text-dark-400 mt-1 flex items-center gap-1.5">
                    <Clock size={14} />
                    {new Date(job.created_at).toLocaleString(i18n.language)}
                  </p>
                </div>
              </div>

              <div className="flex flex-wrap items-center gap-4 sm:gap-8">
                <div className="flex items-center gap-2 min-w-[140px]">
                  {getStatusIcon(job.status)}
                  <span className={cn(
                    "text-sm font-medium",
                    job.status === 'completed' ? "text-green-500" : 
                    job.status === 'failed' ? "text-red-500" : "text-dark-300"
                  )}>
                    {getStatusLabel(job.status)}
                  </span>
                </div>

                {job.status === 'completed' && job.output_files && job.output_files.length > 0 && (
                  <div className="flex flex-col gap-2">
                    {job.output_files.map((file, idx) => (
                      <button
                        key={idx}
                        onClick={() => handleDownload(job.job_id, file)}
                        className="flex items-center gap-2 bg-dark-800 hover:bg-dark-700 text-dark-50 px-4 py-2 rounded-lg text-sm font-semibold transition-all"
                      >
                        <Download size={16} />
                        {job.output_files && job.output_files.length > 1 
                          ? `${t('dashboard_page.download')} ${idx + 1}` 
                          : t('dashboard_page.download')}
                      </button>
                    ))}
                  </div>
                )}
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
};

export default Dashboard;
