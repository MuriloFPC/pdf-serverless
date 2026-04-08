import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import api from '../services/api';
import { 
  Upload, 
  FileText, 
  X, 
  Loader2, 
  ArrowRight, 
  Lock, 
  Unlock, 
  Scissors, 
  Combine,
  AlertCircle,
  Zap
} from 'lucide-react';
import { cn } from '../lib/utils';
import axios from 'axios';
import { useTranslation } from 'react-i18next';

type ProcessType = 'merge' | 'split' | 'protect' | 'unprotect' | 'optimize';

const NewJob: React.FC = () => {
  const { t } = useTranslation();
  const [type, setType] = useState<ProcessType>('merge');
  const [files, setFiles] = useState<File[]>([]);
  const [password, setPassword] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [uploadProgress, setUploadProgress] = useState(0);
  const navigate = useNavigate();

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files) {
      const newFiles = Array.from(e.target.files);
      setFiles(prev => [...prev, ...newFiles]);
    }
  };

  const removeFile = (index: number) => {
    setFiles(prev => prev.filter((_, i) => i !== index));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (files.length === 0) {
      setError(t('new_job_page.error_no_files'));
      return;
    }

    setLoading(true);
    setError('');

    try {
      // 1. Create Job
      const jobResponse = await api.post('/pdf/process', {
        type,
        password,
        ttl: '24h',
        metadata: {}
      });

      const { job_id } = jobResponse.data;

      // 2. Upload files
      let completed = 0;
      for (const file of files) {
        // Get presigned URL
        const urlResponse = await api.get(`/pdf/presigned-url/${job_id}`, {
          params: { filename: file.name }
        });

        const { url } = urlResponse.data;

        // Upload to S3
        await axios.put(url, file, {
          headers: { 'Content-Type': 'application/pdf' },
          onUploadProgress: (progressEvent) => {
             // Basic progress calculation
             const percent = Math.round((progressEvent.loaded * 100) / (progressEvent.total || 1));
             setUploadProgress(Math.round(((completed * 100) + percent) / files.length));
          }
        });
        completed++;
      }

      // 3. Complete Upload
      await api.post(`/pdf/complete-upload/${job_id}`);

      navigate('/');
    } catch (err: any) {
      setError(err.response?.data?.error || t('new_job_page.error_default'));
    } finally {
      setLoading(false);
    }
  };

  const types = [
    { id: 'merge', label: t('new_job_page.tools.merge.label'), icon: Combine, desc: t('new_job_page.tools.merge.desc') },
    { id: 'split', label: t('new_job_page.tools.split.label'), icon: Scissors, desc: t('new_job_page.tools.split.desc') },
    { id: 'protect', label: t('new_job_page.tools.protect.label'), icon: Lock, desc: t('new_job_page.tools.protect.desc') },
    { id: 'unprotect', label: t('new_job_page.tools.unprotect.label'), icon: Unlock, desc: t('new_job_page.tools.unprotect.desc') },
    { id: 'optimize', label: t('new_job_page.tools.optimize.label'), icon: Zap, desc: t('new_job_page.tools.optimize.desc') },
  ];

  return (
    <div className="container mx-auto px-4 py-8 max-w-4xl">
      <div className="mb-8 text-center">
        <h1 className="text-3xl font-bold">{t('new_job_page.title')}</h1>
        <p className="text-dark-400 mt-2">{t('new_job_page.subtitle')}</p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-5 gap-4 mb-8">
        {types.map((t) => (
          <button
            key={t.id}
            onClick={() => setType(t.id as ProcessType)}
            className={cn(
              "p-4 rounded-xl border text-left transition-all flex flex-col gap-3",
              type === t.id 
                ? "bg-blue-600/10 border-blue-500 ring-1 ring-blue-500 text-blue-500" 
                : "bg-dark-900 border-dark-800 text-dark-300 hover:border-dark-700"
            )}
          >
            <t.icon size={24} />
            <div>
              <div className="font-bold">{t.label}</div>
              <div className="text-xs opacity-70 leading-tight mt-1">{t.desc}</div>
            </div>
          </button>
        ))}
      </div>

      <div className="bg-dark-900 border border-dark-800 rounded-2xl p-6 sm:p-8">
        <form onSubmit={handleSubmit} className="space-y-6">
          {error && (
            <div className="bg-red-500/10 border border-red-500/20 text-red-500 p-4 rounded-lg flex items-center gap-3 text-sm">
              <AlertCircle size={18} />
              <span>{error}</span>
            </div>
          )}

          <div>
            <label className="block text-sm font-medium text-dark-300 mb-4">
              {t('new_job_page.tools.merge.label')}
            </label>
            
            <div className="border-2 border-dashed border-dark-800 rounded-xl p-8 text-center hover:border-blue-500/50 transition-colors relative group">
              <input
                type="file"
                multiple
                accept=".pdf"
                onChange={handleFileChange}
                className="absolute inset-0 w-full h-full opacity-0 cursor-pointer"
              />
              <div className="flex flex-col items-center gap-3">
                <div className="h-12 w-12 bg-dark-800 rounded-full flex items-center justify-center text-dark-400 group-hover:text-blue-500 group-hover:bg-blue-500/10 transition-all">
                  <Upload size={24} />
                </div>
                <div className="text-dark-300">
                  <span className="font-bold text-blue-500">{t('new_job_page.select_files')}</span> {t('new_job_page.drag_drop')}
                </div>
                <div className="text-xs text-dark-500">Apenas arquivos .PDF</div>
              </div>
            </div>

            {files.length > 0 && (
              <div className="mt-6 space-y-2">
                {files.map((file, index) => (
                  <div key={index} className="bg-dark-950 border border-dark-800 rounded-lg p-3 flex items-center justify-between group">
                    <div className="flex items-center gap-3 truncate">
                      <FileText size={18} className="text-blue-500 flex-shrink-0" />
                      <span className="text-sm truncate">{file.name}</span>
                      <span className="text-xs text-dark-500">{(file.size / 1024 / 1024).toFixed(2)} MB</span>
                    </div>
                    <button 
                      type="button"
                      onClick={() => removeFile(index)}
                      className="p-1 hover:bg-red-500/10 hover:text-red-500 rounded transition-colors"
                    >
                      <X size={16} />
                    </button>
                  </div>
                ))}
              </div>
            )}
          </div>

          {(type === 'protect' || type === 'unprotect') && (
            <div>
              <label htmlFor="password" className="block text-sm font-medium text-dark-300 mb-1.5">
                {t('new_job_page.password_label')}
              </label>
              <input
                id="password"
                type="password"
                required
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                placeholder="••••••••"
                className="w-full bg-dark-950 border border-dark-800 rounded-lg px-4 py-2.5 focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 outline-none transition-all"
              />
            </div>
          )}

          <div className="pt-4">
            <button
              type="submit"
              disabled={loading || files.length === 0}
              className="w-full bg-blue-600 hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed text-white font-bold py-4 rounded-xl flex items-center justify-center gap-2 transition-all shadow-lg shadow-blue-500/20"
            >
              {loading ? (
                <>
                  <Loader2 className="animate-spin" size={20} />
                  <span>{t('new_job_page.processing')} {uploadProgress}%</span>
                </>
              ) : (
                <>
                  <span>{t('new_job_page.process_button')}</span>
                  <ArrowRight size={20} />
                </>
              )}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

export default NewJob;
