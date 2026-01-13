-- +goose Up
-- +goose StatementBegin
CREATE TABLE public.customer_update_profile (
	id bigserial NOT NULL,
	created timestamp DEFAULT now() NULL,
	customer_name varchar not null,
	customer_email varchar not null,
	reference_number varchar DEFAULT ''::character varying NULL,
	old_data_mobile text NULL,
	new_data_mobile text NULL,
	list_data_changes text DEFAULT '-'::text NULL,
	CONSTRAINT customer_update_profile_pk PRIMARY KEY (id)
);

INSERT INTO public.customer_update_profile (
    created,
    customer_name,
    customer_email,
    reference_number,
    old_data_mobile,
    new_data_mobile,
    list_data_changes
)
SELECT
    now(),
    'Customer ' || gs,
    'customer' || gs || '@example.com',
    'REF' || lpad(gs::text, 6, '0'),
    '{"cellularNo":"628988273336","email":"satoukazumaforwork@gmail.com","emergencyPhone":"62812345678","emergencyContact":"KAZUMA SATOU","jobTitle":"08","jobDesc":"STAFF|STAFF","occupationCode":"003","occupationDesc":"PEGAWAI NEGERI | GOVERNMENT EMPLOYEES","companyName":"ALAM MAKMUR","companyAddr":"Jalan Soetta Ubah Data nya Dong","companyType":"0053","companyDesc":"Perusahaan IT | IT Company","incomeCode":"7","incomeDesc":"Rp 500 juta s/d Rp 1 Milyar | IDR 500 million to IDR 1 billion","addressType":"DOMICILE","address":"JALAN DOMISILI Update","rt":"001","rw":"005","kecamatan":"Pondoksalam","kelurahan":"Desa Tanjung Sari","cityCode":"0103","cityName":"Kab. Purwakarta","provinceCode":"006","provinceName":"Jawa Barat","postalCode":"41115","npwpNumber":"531622336000000","npwpImage":"","kodeKecamatan":"3214091","kodeKelurahan":"3214091006"}',
    '{"cellularNo":"628988273336","email":"satoukazumaforwork@gmail.com","emergencyPhone":"62812345678","emergencyContact":"KAZUMA SATOU","jobTitle":"08","jobDesc":"STAFF|STAFF","occupationCode":"003","occupationDesc":"PEGAWAI NEGERI | GOVERNMENT EMPLOYEES","companyName":"ALAM MAKMUR","companyAddr":"Jalan Soetta Ubah Data nya Dong","companyType":"0053","companyDesc":"Perusahaan IT | IT Company","incomeCode":"7","incomeDesc":"Rp 500 juta s/d Rp 1 Milyar | IDR 500 million to IDR 1 billion","addressType":"DOMICILE","address":"JALAN DOMISILI Update","rt":"001","rw":"005","kecamatan":"Pondoksalam","kelurahan":"Desa Tanjung Sari","cityCode":"0103","cityName":"Kab. Purwakarta","provinceCode":"006","provinceName":"Jawa Barat","postalCode":"41115","npwpNumber":"531622336000000","npwpImage":"","kodeKecamatan":"3214091","kodeKelurahan":"3214091006"}',
    'Alamat Domisili, Alamat Perusahaan'
FROM generate_series(1, 1000000) AS gs;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.customer_update_profile;
-- +goose StatementEnd
