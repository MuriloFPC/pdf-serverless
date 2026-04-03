import React from 'react';
import { Link } from 'react-router-dom';
import { Shield, Lock, CheckCircle2, Database, Cloud, Code, Terminal, ExternalLink, Zap, TrendingDown, RefreshCcw, ShieldCheck, Cpu } from 'lucide-react';
import { useTranslation } from 'react-i18next';

const Landing: React.FC = () => {
  const { t } = useTranslation();

  return (
    <div className="flex flex-col min-h-screen bg-dark-950 text-dark-50">
      {/* Hero Section */}
      <section className="relative py-20 overflow-hidden border-b border-dark-800">
        <div className="container mx-auto px-6">
          <div className="max-w-4xl mx-auto text-center">
            <h1 className="text-5xl md:text-7xl font-extrabold tracking-tight mb-6 bg-gradient-to-r from-blue-400 to-indigo-500 bg-clip-text text-transparent">
              {t('landing.hero.title')}
            </h1>
            <p className="text-xl md:text-2xl text-dark-400 mb-10 leading-relaxed">
              {t('landing.hero.subtitle')}
            </p>
            <div className="flex flex-col sm:flex-row items-center justify-center gap-4">
              <Link
                to="/register"
                className="w-full sm:w-auto px-8 py-4 bg-blue-600 hover:bg-blue-700 text-white rounded-xl font-bold transition-all shadow-lg shadow-blue-500/20 text-center"
              >
                {t('landing.hero.cta_start')}
              </Link>
              <Link
                to="/login"
                className="w-full sm:w-auto px-8 py-4 bg-dark-800 hover:bg-dark-700 text-white border border-dark-700 rounded-xl font-bold transition-all text-center"
              >
                {t('landing.hero.cta_login')}
              </Link>
            </div>
          </div>
        </div>
        
        {/* Background elements */}
        <div className="absolute top-0 left-1/2 -translate-x-1/2 w-[1000px] h-[600px] bg-blue-500/5 blur-[120px] rounded-full -z-10" />
      </section>

      {/* Security Pillars */}
      <section id="features" className="py-24 border-b border-dark-800">
        <div className="container mx-auto px-6">
          <div className="text-center mb-16">
            <h2 className="text-3xl md:text-4xl font-bold mb-4">{t('landing.security.title')}</h2>
            <p className="text-dark-400 max-w-2xl mx-auto">
              {t('landing.security.subtitle')}
            </p>
          </div>
          
          <div className="grid md:grid-cols-3 gap-12">
            <div className="p-8 rounded-2xl bg-dark-900/50 border border-dark-800 hover:border-blue-500/30 transition-colors">
              <div className="w-14 h-14 bg-blue-500/10 rounded-xl flex items-center justify-center mb-6">
                <Shield className="text-blue-500" size={32} />
              </div>
              <h3 className="text-xl font-bold mb-4">{t('landing.security.storage.title')}</h3>
              <p className="text-dark-400 leading-relaxed">
                {t('landing.security.storage.desc')}
              </p>
            </div>
            
            <div className="p-8 rounded-2xl bg-dark-900/50 border border-dark-800 hover:border-blue-500/30 transition-colors">
              <div className="w-14 h-14 bg-indigo-500/10 rounded-xl flex items-center justify-center mb-6">
                <Lock className="text-indigo-500" size={32} />
              </div>
              <h3 className="text-xl font-bold mb-4">{t('landing.security.encryption.title')}</h3>
              <p className="text-dark-400 leading-relaxed">
                {t('landing.security.encryption.desc')}
              </p>
            </div>
            
            <div className="p-8 rounded-2xl bg-dark-900/50 border border-dark-800 hover:border-blue-500/30 transition-colors">
              <div className="w-14 h-14 bg-emerald-500/10 rounded-xl flex items-center justify-center mb-6">
                <Database className="text-emerald-500" size={32} />
              </div>
              <h3 className="text-xl font-bold mb-4">{t('landing.security.governance.title')}</h3>
              <p className="text-dark-400 leading-relaxed">
                {t('landing.security.governance.desc')}
              </p>
            </div>
          </div>
        </div>
      </section>

      {/* How it works */}
      <section className="py-24 bg-dark-900/30 border-b border-dark-800">
        <div className="container mx-auto px-6">
          <div className="grid lg:grid-cols-2 gap-16 items-center">
            <div>
              <h2 className="text-3xl md:text-4xl font-bold mb-8">{t('landing.workflow.title')}</h2>
              <div className="space-y-8">
                <div className="flex gap-6">
                  <div className="flex-shrink-0 w-10 h-10 rounded-full bg-blue-600 flex items-center justify-center font-bold">1</div>
                  <div>
                    <h4 className="text-lg font-bold mb-2">{t('landing.workflow.step1.title')}</h4>
                    <p className="text-dark-400">{t('landing.workflow.step1.desc')}</p>
                  </div>
                </div>
                <div className="flex gap-6">
                  <div className="flex-shrink-0 w-10 h-10 rounded-full bg-blue-600 flex items-center justify-center font-bold">2</div>
                  <div>
                    <h4 className="text-lg font-bold mb-2">{t('landing.workflow.step2.title')}</h4>
                    <p className="text-dark-400">{t('landing.workflow.step2.desc')}</p>
                  </div>
                </div>
                <div className="flex gap-6">
                  <div className="flex-shrink-0 w-10 h-10 rounded-full bg-blue-600 flex items-center justify-center font-bold">3</div>
                  <div>
                    <h4 className="text-lg font-bold mb-2">{t('landing.workflow.step3.title')}</h4>
                    <p className="text-dark-400">{t('landing.workflow.step3.desc')}</p>
                  </div>
                </div>
              </div>
            </div>
            <div className="relative">
              <div className="aspect-video bg-dark-900 rounded-2xl border border-dark-700 shadow-2xl flex items-center justify-center p-8 overflow-hidden group">
                <div className="grid grid-cols-2 gap-4 w-full">
                  <div className="h-20 bg-dark-800 rounded-lg animate-pulse" />
                  <div className="h-20 bg-dark-800 rounded-lg animate-pulse delay-75" />
                  <div className="h-20 bg-dark-800 rounded-lg animate-pulse delay-150" />
                  <div className="h-20 bg-dark-800 rounded-lg animate-pulse delay-200" />
                </div>
                <div className="absolute inset-0 bg-gradient-to-t from-dark-950/80 to-transparent flex items-end justify-center pb-8">
                  <div className="px-6 py-2 bg-blue-600 rounded-full text-sm font-bold shadow-lg">{t('landing.workflow.realtime_dashboard')}</div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* Features List */}
      <section className="py-24 border-b border-dark-800">
        <div className="container mx-auto px-6">
          <div className="grid md:grid-cols-2 lg:grid-cols-4 gap-8">
            <div className="flex items-start gap-3">
              <CheckCircle2 className="text-blue-500 mt-1" size={20} />
              <span>{t('landing.features.merge')}</span>
            </div>
            <div className="flex items-start gap-3">
              <CheckCircle2 className="text-blue-500 mt-1" size={20} />
              <span>{t('landing.features.split')}</span>
            </div>
            <div className="flex items-start gap-3">
              <CheckCircle2 className="text-blue-500 mt-1" size={20} />
              <span>{t('landing.features.protect')}</span>
            </div>
            <div className="flex items-start gap-3">
              <CheckCircle2 className="text-blue-500 mt-1" size={20} />
              <span>{t('landing.features.unprotect')}</span>
            </div>
          </div>
        </div>
      </section>

      {/* Comparison Section */}
      <section className="py-24 border-b border-dark-800 bg-dark-900/10">
        <div className="container mx-auto px-6">
          <div className="text-center mb-16">
            <h2 className="text-3xl md:text-4xl font-bold mb-4">{t('landing.comparison.title')}</h2>
            <p className="text-dark-400 max-w-2xl mx-auto">
              {t('landing.comparison.subtitle')}
            </p>
          </div>

          <div className="grid lg:grid-cols-2 gap-8 max-w-5xl mx-auto">
            {/* Platform (SaaS) */}
            <div className="relative p-8 rounded-3xl bg-blue-600/5 border-2 border-blue-500/30 overflow-hidden group hover:border-blue-500 transition-all">
              <div className="absolute top-0 right-0 p-4">
                <span className="px-3 py-1 bg-blue-600 text-xs font-bold rounded-full uppercase tracking-wider">{t('landing.comparison.platform.tag')}</span>
              </div>
              
              <div className="flex items-center gap-4 mb-8">
                <div className="w-12 h-12 bg-blue-500 rounded-2xl flex items-center justify-center shadow-lg shadow-blue-500/20">
                  <Zap className="text-white" size={24} />
                </div>
                <div>
                  <h3 className="text-2xl font-bold">{t('landing.comparison.platform.title')}</h3>
                  <p className="text-blue-400 text-sm">{t('landing.comparison.platform.subtitle')}</p>
                </div>
              </div>

              <ul className="space-y-6 mb-10">
                <li className="flex gap-4">
                  <TrendingDown className="text-blue-500 shrink-0" size={24} />
                  <div>
                    <h5 className="font-bold">{t('landing.comparison.platform.item1.title')}</h5>
                    <p className="text-sm text-dark-400">{t('landing.comparison.platform.item1.desc')}</p>
                  </div>
                </li>
                <li className="flex gap-4">
                  <RefreshCcw className="text-blue-500 shrink-0" size={24} />
                  <div>
                    <h5 className="font-bold">{t('landing.comparison.platform.item2.title')}</h5>
                    <p className="text-sm text-dark-400">{t('landing.comparison.platform.item2.desc')}</p>
                  </div>
                </li>
                <li className="flex gap-4">
                  <ShieldCheck className="text-blue-500 shrink-0" size={24} />
                  <div>
                    <h5 className="font-bold">{t('landing.comparison.platform.item3.title')}</h5>
                    <p className="text-sm text-dark-400">{t('landing.comparison.platform.item3.desc')}</p>
                  </div>
                </li>
              </ul>
              
              <Link
                to="/register"
                className="block w-full py-4 bg-blue-600 hover:bg-blue-700 text-white rounded-xl font-bold text-center transition-all"
              >
                {t('landing.comparison.platform.cta')}
              </Link>
            </div>

            {/* Self-Hosted */}
            <div className="p-8 rounded-3xl bg-dark-900 border border-dark-800 hover:border-dark-600 transition-all">
              <div className="flex items-center gap-4 mb-8">
                <div className="w-12 h-12 bg-dark-800 rounded-2xl flex items-center justify-center border border-dark-700">
                  <Terminal className="text-dark-300" size={24} />
                </div>
                <div>
                  <h3 className="text-2xl font-bold text-dark-100">{t('landing.comparison.self_hosted.title')}</h3>
                  <p className="text-dark-500 text-sm">{t('landing.comparison.self_hosted.subtitle')}</p>
                </div>
              </div>

              <ul className="space-y-6 mb-10">
                <li className="flex gap-4 opacity-70">
                  <CheckCircle2 className="text-dark-400 shrink-0" size={24} />
                  <div>
                    <h5 className="font-bold">{t('landing.comparison.self_hosted.item1.title')}</h5>
                    <p className="text-sm text-dark-500">{t('landing.comparison.self_hosted.item1.desc')}</p>
                  </div>
                </li>
                <li className="flex gap-4 opacity-70">
                  <CheckCircle2 className="text-dark-400 shrink-0" size={24} />
                  <div>
                    <h5 className="font-bold">{t('landing.comparison.self_hosted.item2.title')}</h5>
                    <p className="text-sm text-dark-500">{t('landing.comparison.self_hosted.item2.desc')}</p>
                  </div>
                </li>
                <li className="flex gap-4 opacity-70">
                  <CheckCircle2 className="text-dark-400 shrink-0" size={24} />
                  <div>
                    <h5 className="font-bold">{t('landing.comparison.self_hosted.item3.title')}</h5>
                    <p className="text-sm text-dark-500">{t('landing.comparison.self_hosted.item3.desc')}</p>
                  </div>
                </li>
              </ul>

              <a 
                href="https://github.com/MuriloFPC/pdf-serverless" 
                target="_blank" 
                rel="noopener noreferrer"
                className="block w-full py-4 bg-dark-800 hover:bg-dark-700 text-white border border-dark-700 rounded-xl font-bold text-center transition-all"
              >
                {t('landing.comparison.self_hosted.cta')}
              </a>
            </div>
          </div>
        </div>
      </section>

      {/* API First Section */}
      <section className="py-24 border-b border-dark-800 bg-dark-900/10">
        <div className="container mx-auto px-6">
          <div className="flex flex-col lg:flex-row items-center gap-16">
            <div className="lg:w-1/2">
              <div className="inline-flex items-center gap-2 px-4 py-2 bg-emerald-500/10 text-emerald-400 rounded-full text-sm font-medium mb-6 border border-emerald-500/20">
                <Cpu size={16} />
                <span>{t('landing.api_first.badge')}</span>
              </div>
              <h2 className="text-3xl md:text-4xl font-bold mb-6">{t('landing.api_first.title')}</h2>
              <p className="text-xl text-dark-400 mb-8 leading-relaxed">
                {t('landing.api_first.subtitle')}
              </p>
              
              <ul className="space-y-4">
                <li className="flex items-start gap-3">
                  <div className="mt-1 bg-emerald-500/20 p-1 rounded">
                    <CheckCircle2 className="text-emerald-500" size={16} />
                  </div>
                  <div>
                    <span className="font-bold block text-dark-100">{t('landing.api_first.item1.title')}</span>
                    <p className="text-dark-400 text-sm">{t('landing.api_first.item1.desc')}</p>
                  </div>
                </li>
                <li className="flex items-start gap-3">
                  <div className="mt-1 bg-emerald-500/20 p-1 rounded">
                    <CheckCircle2 className="text-emerald-500" size={16} />
                  </div>
                  <div>
                    <span className="font-bold block text-dark-100">{t('landing.api_first.item2.title')}</span>
                    <p className="text-dark-400 text-sm">{t('landing.api_first.item2.desc')}</p>
                  </div>
                </li>
                <li className="flex items-start gap-3">
                  <div className="mt-1 bg-emerald-500/20 p-1 rounded">
                    <CheckCircle2 className="text-emerald-500" size={16} />
                  </div>
                  <div>
                    <span className="font-bold block text-dark-100">{t('landing.api_first.item3.title')}</span>
                    <p className="text-dark-400 text-sm">{t('landing.api_first.item3.desc')}</p>
                  </div>
                </li>
              </ul>
            </div>
            
            <div className="lg:w-1/2 w-full">
              <div className="bg-dark-900 rounded-2xl border border-dark-700 p-6 font-mono text-sm overflow-hidden shadow-2xl relative">
                <div className="flex gap-2 mb-4 border-b border-dark-800 pb-4">
                  <div className="w-3 h-3 rounded-full bg-red-500/50"></div>
                  <div className="w-3 h-3 rounded-full bg-yellow-500/50"></div>
                  <div className="w-3 h-3 rounded-full bg-green-500/50"></div>
                  <span className="ml-2 text-dark-600">POST /api/jobs</span>
                </div>
                <div className="text-blue-400">curl</div> <span className="text-dark-300">-X POST</span> https://api.pdf-serverless.com/jobs \<br/>
                &nbsp;&nbsp;<span className="text-dark-300">-H</span> <span className="text-emerald-400">"Authorization: Bearer $TOKEN"</span> \<br/>
                &nbsp;&nbsp;<span className="text-dark-300">-F</span> <span className="text-emerald-400">"type=merge"</span> \<br/>
                &nbsp;&nbsp;<span className="text-dark-300">-F</span> <span className="text-emerald-400">"files=@document1.pdf"</span> \<br/>
                &nbsp;&nbsp;<span className="text-dark-300">-F</span> <span className="text-emerald-400">"files=@document2.pdf"</span><br/>
                <br/>
                <span className="text-dark-600">// {t('landing.api_first.response_json')}</span><br/>
                <span className="text-purple-400">{"{"}</span><br/>
                &nbsp;&nbsp;<span className="text-blue-300">"id"</span>: <span className="text-emerald-400">"job_82k91..."</span>,<br/>
                &nbsp;&nbsp;<span className="text-blue-300">"status"</span>: <span className="text-emerald-400">"processing"</span><br/>
                <span className="text-purple-400">{"}"}</span>
                
                {/* Decorative glow */}
                <div className="absolute -bottom-10 -right-10 w-40 h-40 bg-blue-500/10 blur-3xl rounded-full"></div>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* Open Source / Self-Host Section */}
      <section className="py-24 bg-gradient-to-b from-dark-900/50 to-dark-950">
        <div className="container mx-auto px-6 text-center">
          <div className="max-w-3xl mx-auto">
            <div className="inline-flex items-center gap-2 px-4 py-2 bg-blue-500/10 text-blue-400 rounded-full text-sm font-medium mb-8 border border-blue-500/20">
              <Code size={16} />
              <span>{t('landing.open_source.badge')}</span>
            </div>
            <h2 className="text-3xl md:text-5xl font-bold mb-6">{t('landing.open_source.title')}</h2>
            <p className="text-xl text-dark-400 mb-10 leading-relaxed">
              {t('landing.open_source.subtitle')}
            </p>
            <div className="grid sm:grid-cols-2 gap-6 mb-12">
              <div className="p-6 rounded-2xl bg-dark-900 border border-dark-800 text-left">
                <Cloud className="text-blue-500 mb-4" size={24} />
                <h4 className="font-bold mb-2">{t('landing.open_source.item1.title')}</h4>
                <p className="text-sm text-dark-500">{t('landing.open_source.item1.desc')}</p>
              </div>
              <div className="p-6 rounded-2xl bg-dark-900 border border-dark-800 text-left">
                <Terminal className="text-blue-500 mb-4" size={24} />
                <h4 className="font-bold mb-2">{t('landing.open_source.item2.title')}</h4>
                <p className="text-sm text-dark-500">{t('landing.open_source.item2.desc')}</p>
              </div>
            </div>
            <a 
              href="https://github.com/MuriloFPC/pdf-serverless" 
              target="_blank" 
              rel="noopener noreferrer"
              className="inline-flex items-center gap-3 px-8 py-4 bg-white text-dark-950 hover:bg-dark-200 rounded-xl font-bold transition-all shadow-xl"
            >
              <ExternalLink size={20} />
              {t('landing.open_source.cta')}
            </a>
          </div>
        </div>
      </section>

      {/* Footer */}
      <footer className="py-12 border-t border-dark-800">
        <div className="container mx-auto px-6 text-center text-dark-500 text-sm">
          <p>{t('landing.footer.text')}</p>
          <p className="mt-2 italic">{t('landing.footer.subtext')}</p>
        </div>
      </footer>
    </div>
  );
};

export default Landing;
