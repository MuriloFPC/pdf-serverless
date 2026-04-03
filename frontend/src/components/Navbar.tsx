import React, { useState, useRef, useEffect } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { FileText, LogOut, LayoutDashboard, PlusCircle, Globe, ChevronDown } from 'lucide-react';
import { useTranslation } from 'react-i18next';

const Navbar: React.FC = () => {
  const { user, logout } = useAuth();
  const navigate = useNavigate();
  const { t, i18n } = useTranslation();
  const [langOpen, setLangOpen] = useState(false);
  const dropdownRef = useRef<HTMLDivElement>(null);

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  const changeLanguage = (lng: string) => {
    i18n.changeLanguage(lng);
    setLangOpen(false);
  };

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target as Node)) {
        setLangOpen(false);
      }
    };
    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, []);

  const languages = [
    { code: 'en', name: 'English', flag: '🇺🇸' },
    { code: 'pt', name: 'Português', flag: '🇧🇷' },
    { code: 'es', name: 'Español', flag: '🇪🇸' },
  ];

  const currentLang = languages.find(l => l.code === i18n.language.split('-')[0]) || languages[0];

  return (
    <nav className="border-b border-dark-800 bg-dark-900/50 backdrop-blur-md sticky top-0 z-50">
      <div className="container mx-auto px-4 h-16 flex items-center justify-between">
        <Link to="/" className="flex items-center gap-2 font-bold text-xl text-blue-500">
          <FileText size={24} />
          <span>PDF Serverless</span>
        </Link>

        <div className="flex items-center gap-4 sm:gap-6">
          {/* Language Selector */}
          <div className="relative" ref={dropdownRef}>
            <button
              onClick={() => setLangOpen(!langOpen)}
              className="flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-dark-300 hover:text-dark-50 hover:bg-dark-800 transition-all text-sm font-medium"
            >
              <Globe size={16} />
              <span className="hidden sm:inline">{currentLang.name}</span>
              <span className="sm:hidden">{currentLang.code.toUpperCase()}</span>
              <ChevronDown size={14} className={`transition-transform ${langOpen ? 'rotate-180' : ''}`} />
            </button>

            {langOpen && (
              <div className="absolute right-0 mt-2 w-40 bg-dark-900 border border-dark-800 rounded-xl shadow-2xl py-2 z-50 animate-in fade-in zoom-in duration-200">
                {languages.map((lang) => (
                  <button
                    key={lang.code}
                    onClick={() => changeLanguage(lang.code)}
                    className={`w-full flex items-center gap-3 px-4 py-2 text-sm transition-colors ${
                      i18n.language.startsWith(lang.code)
                        ? 'text-blue-500 bg-blue-500/5 font-bold'
                        : 'text-dark-300 hover:bg-dark-800'
                    }`}
                  >
                    <span>{lang.flag}</span>
                    <span>{lang.name}</span>
                  </button>
                ))}
              </div>
            )}
          </div>

          {user ? (
            <div className="flex items-center gap-4 sm:gap-6">
              <Link to="/" className="flex items-center gap-1.5 text-dark-300 hover:text-dark-50 transition-colors">
                <LayoutDashboard size={18} />
                <span className="hidden md:inline">{t('nav.dashboard')}</span>
              </Link>
              <Link to="/new-job" className="flex items-center gap-1.5 text-dark-300 hover:text-dark-50 transition-colors">
                <PlusCircle size={18} />
                <span className="hidden md:inline">{t('nav.new_job')}</span>
              </Link>
              
              <div className="h-6 w-px bg-dark-800 mx-1 sm:mx-2" />
              
              <div className="flex items-center gap-3">
                <span className="text-sm text-dark-400 hidden lg:inline">{user.email}</span>
                <button
                  onClick={handleLogout}
                  className="p-2 text-dark-400 hover:text-red-400 hover:bg-red-400/10 rounded-lg transition-all"
                  title={t('common.logout')}
                >
                  <LogOut size={20} />
                </button>
              </div>
            </div>
          ) : (
            <div className="flex items-center gap-2 sm:gap-4">
              <Link 
                to="/login" 
                className="px-3 sm:px-4 py-2 text-dark-300 hover:text-dark-50 transition-colors font-medium text-sm sm:text-base"
              >
                {t('nav.sign_in')}
              </Link>
              <Link 
                to="/register" 
                className="px-4 sm:px-5 py-2 sm:py-2.5 bg-blue-600 hover:bg-blue-700 text-white rounded-lg font-bold transition-all shadow-lg shadow-blue-500/20 text-sm sm:text-base whitespace-nowrap"
              >
                {t('nav.get_started')}
              </Link>
            </div>
          )}
        </div>
      </div>
    </nav>
  );
};

export default Navbar;
