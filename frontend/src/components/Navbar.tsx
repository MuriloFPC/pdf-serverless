import React from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { FileText, LogOut, LayoutDashboard, PlusCircle } from 'lucide-react';

const Navbar: React.FC = () => {
  const { user, logout } = useAuth();
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  return (
    <nav className="border-b border-dark-800 bg-dark-900/50 backdrop-blur-md sticky top-0 z-50">
      <div className="container mx-auto px-4 h-16 flex items-center justify-between">
        <Link to="/" className="flex items-center gap-2 font-bold text-xl text-blue-500">
          <FileText size={24} />
          <span>PDF Serverless</span>
        </Link>

        {user && (
          <div className="flex items-center gap-6">
            <Link to="/" className="flex items-center gap-1.5 text-dark-300 hover:text-dark-50 transition-colors">
              <LayoutDashboard size={18} />
              <span>Dashboard</span>
            </Link>
            <Link to="/new-job" className="flex items-center gap-1.5 text-dark-300 hover:text-dark-50 transition-colors">
              <PlusCircle size={18} />
              <span>Novo Processo</span>
            </Link>
            
            <div className="h-6 w-px bg-dark-800 mx-2" />
            
            <div className="flex items-center gap-4">
              <span className="text-sm text-dark-400 hidden sm:inline">{user.email}</span>
              <button
                onClick={handleLogout}
                className="p-2 text-dark-400 hover:text-red-400 hover:bg-red-400/10 rounded-lg transition-all"
                title="Sair"
              >
                <LogOut size={20} />
              </button>
            </div>
          </div>
        )}
      </div>
    </nav>
  );
};

export default Navbar;
