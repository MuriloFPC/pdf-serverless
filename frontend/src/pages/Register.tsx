import React, { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import api from '../services/api';
import { UserPlus, LogIn, AlertCircle, Loader2, CheckCircle } from 'lucide-react';

const Register: React.FC = () => {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [error, setError] = useState('');
  const [success, setSuccess] = useState(false);
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    
    if (password !== confirmPassword) {
      setError('As senhas não coincidem');
      return;
    }

    setLoading(true);

    try {
      await api.post('/auth/register', { email, password });
      setSuccess(true);
      setTimeout(() => navigate('/login'), 2000);
    } catch (err: any) {
      setError(err.response?.data?.error || 'Erro ao criar conta. Tente novamente.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-[calc(100vh-64px)] flex items-center justify-center px-4">
      <div className="w-full max-w-md space-y-8 bg-dark-900 p-8 rounded-2xl border border-dark-800 shadow-xl">
        <div className="text-center">
          <div className="mx-auto h-12 w-12 bg-green-500/10 rounded-xl flex items-center justify-center text-green-500 mb-4">
            <UserPlus size={28} />
          </div>
          <h2 className="text-3xl font-bold">Crie sua conta</h2>
          <p className="mt-2 text-dark-400">Comece a processar seus PDFs agora mesmo</p>
        </div>

        {success ? (
          <div className="bg-green-500/10 border border-green-500/20 text-green-500 p-6 rounded-lg flex flex-col items-center gap-3 text-center">
            <CheckCircle size={48} className="mb-2" />
            <h3 className="font-bold text-lg">Conta criada com sucesso!</h3>
            <p className="text-sm">Redirecionando para o login...</p>
          </div>
        ) : (
          <>
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
                    E-mail
                  </label>
                  <input
                    id="email"
                    type="email"
                    required
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                    className="w-full bg-dark-950 border border-dark-800 rounded-lg px-4 py-2.5 focus:ring-2 focus:ring-green-500/20 focus:border-green-500 outline-none transition-all placeholder:text-dark-600"
                    placeholder="seu@email.com"
                  />
                </div>
                <div>
                  <label htmlFor="password" className="block text-sm font-medium text-dark-300 mb-1.5">
                    Senha
                  </label>
                  <input
                    id="password"
                    type="password"
                    required
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    className="w-full bg-dark-950 border border-dark-800 rounded-lg px-4 py-2.5 focus:ring-2 focus:ring-green-500/20 focus:border-green-500 outline-none transition-all placeholder:text-dark-600"
                    placeholder="••••••••"
                  />
                </div>
                <div>
                  <label htmlFor="confirmPassword" className="block text-sm font-medium text-dark-300 mb-1.5">
                    Confirmar Senha
                  </label>
                  <input
                    id="confirmPassword"
                    type="password"
                    required
                    value={confirmPassword}
                    onChange={(e) => setConfirmPassword(e.target.value)}
                    className="w-full bg-dark-950 border border-dark-800 rounded-lg px-4 py-2.5 focus:ring-2 focus:ring-green-500/20 focus:border-green-500 outline-none transition-all placeholder:text-dark-600"
                    placeholder="••••••••"
                  />
                </div>
              </div>

              <button
                type="submit"
                disabled={loading}
                className="w-full bg-green-600 hover:bg-green-700 disabled:opacity-50 disabled:cursor-not-allowed text-white font-semibold py-3 rounded-lg flex items-center justify-center gap-2 transition-all shadow-lg shadow-green-500/20"
              >
                {loading ? <Loader2 className="animate-spin" size={20} /> : 'Criar Conta'}
              </button>
            </form>

            <p className="text-center text-sm text-dark-400">
              Já tem uma conta?{' '}
              <Link to="/login" className="text-green-500 hover:text-green-400 font-medium inline-flex items-center gap-1 transition-colors">
                Faça login <LogIn size={14} />
              </Link>
            </p>
          </>
        )}
      </div>
    </div>
  );
};

export default Register;
