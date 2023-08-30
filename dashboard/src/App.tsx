import { Provider } from 'react-redux';
import {
  createBrowserRouter,
  Link,
  NavLink,
  Outlet,
  RouterProvider,
} from "react-router-dom";

import { CloudEvents } from "./components/CloudEvents/CloudEvents";
import { Workflow } from "./components/Workflow/Workflow";
import { store } from "./store/configure";
import './App.css';


const App = () => {
  const router = createBrowserRouter([
    {
      path: "/", // cloud-events
      element: <Root />,
      errorElement: <Error />,
      children: [
        {
          path: "/cloud-events", // cloud-events
          element: <CloudEvents />,
        },
        {
          path: "/workflows", // cloud-events
          element: <Workflow />,
        }
      ]
    },
  ]);

  return (<>
    <Provider store={store}>
      <RouterProvider router={router} />
    </Provider>
  </>);
}

const Root = () => {
  return (<>
    <header>
      <nav className="navbar navbar-expand-md navbar-dark bg-dark">
        <div className="container-fluid">
          <Link className="navbar-brand" to="/">Fermions</Link>
          <button
            className="navbar-toggler"
            type="button"
            data-bs-toggle="collapse"
            data-bs-target="#navbarCollapse"
            aria-controls="navbarCollapse"
            aria-expanded="false"
            aria-label="Toggle navigation"
          >
            <span className="navbar-toggler-icon"></span>
          </button>
          <div className="collapse navbar-collapse" id="navbarCollapse">
            <ul className="navbar-nav me-auto mb-2 mb-md-0">
            <li className="nav-item">
                <NavLink
                  to="/workflows"
                  className={({ isActive, isPending }) => {
                    return "nav-link " + (isActive ? "active" : "");
                  }}
                >
                  Workflows
                </NavLink>
              </li>
              <li className="nav-item">
                <NavLink
                  className={({ isActive, isPending }) => {
                    return "nav-link " + (isActive ? "active" : "");
                  }}
                  to="/cloud-events"
                >
                  Cloud Events
                </NavLink>
              </li>
            </ul>
          </div>
        </div>
      </nav>
    </header>

    <main style={{height: "calc(100% - 56px - 56px)"}}>
      <Outlet />
    </main>

    <footer className="footer mt-auto py-3 bg-body-tertiary">
      <div className="container">
        <span className="text-body-secondary fw-light">Copyright &#169; 2023 The Fermions Authors</span>
      </div>
    </footer>
  </>);
}

const Error = () => {
  return (
    <div>error</div>
  );
};

export default App;
