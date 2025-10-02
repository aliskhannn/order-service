import { useState } from "react";
import styles from './SearchBar.module.css';

interface SearchBarProps {
	onSearch: (query: string) => void;
}

export const SearchBar: React.FC<SearchBarProps> = ({ onSearch }) => {
	const [query, setQuery] = useState('')
	const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
		setQuery(e.target.value);
	}
	const handleSearch = () => {
		onSearch(query);
	};

	return (
		<div className={styles.searchBar}>
			<input
				type="text" 
				value={query}
				name="search"
				onChange={handleInputChange}
				placeholder="Search order by id"
			/>
			<button onClick={handleSearch}>Search</button>
		</div>
	)
}