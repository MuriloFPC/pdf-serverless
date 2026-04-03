import React, { useState, useEffect } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import api from '../services/api';
import { useAuth } from '../context/AuthContext';
import { LogIn, UserPlus, AlertCircle, Loader2 } from 'lucide-react';
import { useTranslation } from 'react-i18next';

const Login: React.FC = () => {
  const { t } = useTranslation();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const { login, user } = useAuth();
  const navigate = useNavigate();

  useEffect(() => {
    if (user) {
      navigate('/');
    }
  }, [user, navigate]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      console.log('Tentando login para:', email);
      const response = await api.post('/auth/login', { email, password });
      console.log('Resposta completa do login:', response.data);
      
      const { token: newToken, user: newUser } = response.data;
      
      if (!newToken || !newUser) {
        console.error('Login retornou dados incompletos:', { newToken, newUser });
        setError(t('login_page.error_incomplete'));
        return;
      }
      
      console.log('Login bem-sucedido, chamando login context');
      login(newToken, newUser);
      console.log('Navegando para /');
      navigate('/');
    } catch (err: any) {
      console.error('Erro no Login:', err);
      setError(err.response?.data?.error || t('login_page.error_default'));
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-[calc(100vh-64px)] flex items-center justify-center px-4">
      <div className="w-full max-w-md space-y-8 bg-dark-900 p-8 rounded-2xl border border-dark-800 shadow-xl">
        <div className="text-center">
          <div className="mx-auto h-12 w-12 bg-blue-500/10 rounded-xl flex items-center justify-center text-blue-500 mb-4">
            <LogIn size={28} />
          </div>
          <h2 className="text-3xl font-bold">{t('login_page.title')}</h2>
          <p className="mt-2 text-dark-400">{t('login_page.subtitle')}</p>
        </div>

        <form className="mt-8 space-y-6" onSubmit={handleSubmit}>
          {error && (
            <div className="bg-red-500/10 border border-red-500/20 text-red-500 p-4 rounded-lg flex items-center gap-3 text-sm">
              <AlertCircle size={18} />
              <span>{error}</span>
            </div>
          )}

          <div className="space-y-4">
            <div>
              <label htmlFor="email" className="block text-sm font-medium text-dark-300 mb-1.5">
                {t('common.email')}
              </label>
              <input
                id="email"
                type="email"
                required
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                className="w-full bg-dark-950 border border-dark-800 rounded-lg px-4 py-2.5 focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 outline-none transition-all placeholder:text-dark-600"
                placeholder={t('login_page.email_placeholder')}
              />
            </div>
            <div>
              <label htmlFor="password" className="block text-sm font-medium text-dark-300 mb-1.5">
                {t('common.password')}
              </label>
              <input
                id="password"
                type="password"
                required
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                className="w-full bg-dark-950 border border-dark-800 rounded-lg px-4 py-2.5 focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 outline-none transition-all placeholder:text-dark-600"
                placeholder={t('login_page.password_placeholder')}
              />
            </div>
          </div>

          <button
            type="submit"
            disabled={loading}
            className="w-full bg-blue-600 hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed text-white font-semibold py-3 rounded-lg flex items-center justify-center gap-2 transition-all shadow-lg shadow-blue-500/20"
          >
            {loading ? <Loader2 className="animate-spin" size={20} /> : t('common.login')}
          </button>
        </form>

        <p className="text-center text-sm text-dark-400">
          {t('login_page.no_account')}{' '}
          <Link to="/register" className="text-blue-500 hover:text-blue-400 font-medium inline-flex items-center gap-1 transition-colors">
            {t('login_page.signup_link')} <UserPlus size={14} />
          </Link>
        </p>
      </div>
    </div>
  );
};

export default Login;
