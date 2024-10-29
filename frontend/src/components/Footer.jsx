import React from 'react';
import { Link } from 'react-router-dom';
import '../styles/footer.scss';

function Footer() {
    return (
        <footer className="custom-footer py-4">
            <div className="container text-center text-md-start">
                <div className="row footer-cust">
                {/* Copyright and Additional Links */}
                <div className="text-center pt-3">
                    <p>© {new Date().getFullYear()} DocuDefense. All rights reserved.</p>
                    <p>
                        Need help or have questions? <Link to="/contact" className="footer-link">Contact Support</Link>
                    </p>
                </div>
                </div>
            </div>
        </footer>
    );
}

export default Footer;
